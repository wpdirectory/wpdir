package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	retry "github.com/giantswarm/retry-go"
	"github.com/wpdirectory/wpdir/internal/client"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/files"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/searcher"
	"github.com/wpdirectory/wpdir/internal/ulid"
	"github.com/wpdirectory/wpdir/internal/utils"
	"github.com/wpdirectory/wporg"
)

const (
	archiveURL = "http://downloads.wordpress.org/plugin/%s.latest-stable.zip?nostats=1"
)

// Plugin ...
type Plugin struct {
	Slug                   string     `json:"slug"`
	Name                   string     `json:"name,omitempty"`
	Version                string     `json:"version,omitempty"`
	Author                 string     `json:"author,omitempty"`
	AuthorProfile          string     `json:"author_profile,omitempty"`
	Contributors           [][]string `json:"contributors,omitempty"`
	Requires               string     `json:"requires,omitempty"`
	Tested                 string     `json:"tested,omitempty"`
	RequiresPHP            string     `json:"requires_php,omitempty"`
	Rating                 int        `json:"rating,omitempty"`
	Ratings                []Rating   `json:"ratings,omitempty"`
	NumRatings             int        `json:"num_ratings,omitempty"`
	SupportThreads         int        `json:"support_threads,omitempty"`
	SupportThreadsResolved int        `json:"support_threads_resolved,omitempty"`
	ActiveInstalls         int        `json:"active_installs,omitempty"`
	Downloaded             int        `json:"downloaded,omitempty"`
	LastUpdated            string     `json:"last_updated,omitempty"`
	Added                  string     `json:"added,omitempty"`
	Homepage               string     `json:"homepage,omitempty"`
	Sections               struct {
		Description string `json:"description,omitempty"`
		Faq         string `json:"faq,omitempty"`
		Changelog   string `json:"changelog,omitempty"`
		Screenshots string `json:"screenshots,omitempty"`
	} `json:"sections,omitempty"`
	ShortDescription string             `json:"short_description,omitempty"`
	DownloadLink     string             `json:"download_link,omitempty"`
	Screenshots      []Screenshot       `json:"screenshots,omitempty"`
	Tags             [][]string         `json:"tags,omitempty"`
	StableTag        string             `json:"stable_tag,omitempty"`
	Versions         [][]string         `json:"versions,omitempty"`
	DonateLink       string             `json:"donate_link,omitempty"`
	Status           status             `json:"status,omitempty"`
	Searcher         *searcher.Searcher `json:"-"`
	Stats            *files.Stats       `json:"stats,omitempty"`
	indexed          bool
	sync.RWMutex
}

// Rating contains information about ratings of a specific star level (0-5)
type Rating struct {
	Stars  string `json:"stars"`
	Number int    `json:"number"`
}

// Screenshot contains the source and caption of a screenshot
type Screenshot struct {
	Src     string `json:"src"`
	Caption string `json:"caption"`
}

type status int

const (
	// Open shows we have files and API info stored
	Open status = iota
	// Closed shows we cannot get data
	Closed
)

// New returns a new plugin struct.
func New(slug string) *Plugin {

	return &Plugin{
		Slug:   slug,
		Status: Closed,
	}

}

// GetStatus returns the Status as a string
func (p *Plugin) GetStatus() string {
	p.RLock()
	defer p.RUnlock()

	switch p.Status {
	case Open:
		return "Open"
	case Closed:
		return "Closed"
	default:
		return "Invalid Status"
	}
}

// HasIndex returns the index status
func (p *Plugin) HasIndex() bool {
	p.RLock()
	defer p.RUnlock()
	return p.indexed
}

// SetIndexed sets the indexed value
func (p *Plugin) SetIndexed(idx bool) {
	p.Lock()
	defer p.Unlock()
	p.indexed = idx
}

// LoadAPIData updates the Plugin struct with data from an HTTP API
func (p *Plugin) LoadAPIData() error {
	var data []byte
	var err error

	fetch := func() error {
		data, err = p.getAPIData()
		return err
	}

	err = retry.Do(fetch, retry.Timeout(15*time.Second), retry.MaxTries(3), retry.Sleep(5*time.Second))
	if err != nil || data == nil {
		p.Status = Closed
		return err
	}

	err = json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	return nil
}

// GetAPIData ...
func (p *Plugin) getAPIData() ([]byte, error) {
	var data []byte
	var err error

	api := wporg.NewClient()
	data, err = api.GetInfo("plugins", p.Slug)

	return data, err
}

// Update ...
func (p *Plugin) Update() error {
	p.RLock()
	defer p.RUnlock()

	bytes, err := p.getArchive()
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		p.Status = Closed
		return nil
	}

	ref, stats, err := p.processArchive(bytes)
	if err != nil {
		return err
	}

	// Store File Stats
	p.Stats = stats

	if p.Searcher == nil {
		// New Searcher
		sr, err := searcher.New(ref)
		if err != nil {
			return err
		}
		p.Searcher = sr
	} else {
		// Use Existing Searcher
		idx, err := ref.Open()
		if err != nil {
			return err
		}

		err = p.Searcher.SwapIndexes(idx)
		if err != nil {
			return err
		}
	}

	return nil
}

// getArchive ...
func (p *Plugin) getArchive() ([]byte, error) {

	var content []byte
	var err error

	client := client.GetAPI()
	URL := fmt.Sprintf(archiveURL, p.Slug)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println(err)
		return content, err
	}

	// Set User-Agent
	req.Header.Set("User-Agent", "wpdirectory/0.1.0")

	resp, err := client.Do(req)
	if err != nil {
		return content, err
	}
	defer utils.CheckClose(resp.Body, &err)

	if resp.StatusCode != 200 {

		// Code 404 is acceptable, it means the plugin/theme is no longer available.
		if resp.StatusCode == 404 {
			return content, nil
		}

		log.Printf("Downloading the extension '%s' failed. Response code: %d\n", p.Name, resp.StatusCode)

		return content, err

	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}

	return content, nil

}

// processArchive ...
func (p *Plugin) processArchive(archive []byte) (*index.IndexRef, *files.Stats, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	id := ulid.New()
	dst := filepath.Join(wd, "data", "index", "plugins", id)
	opts := &index.IndexOptions{
		ExcludeDotFiles: true,
	}

	ref, stats, err := index.BuildFromZip(opts, archive, dst, p.Slug)
	if err != nil {
		return nil, nil, err
	}

	return ref, stats, nil
}

// Save ...
// TODO: Wrap struct to allow lock during Marshal
func (p *Plugin) Save() error {
	p.RLock()
	defer p.RUnlock()

	bytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	db.PutToBucket(p.Slug, bytes, "plugins")

	return nil
}
