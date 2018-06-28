package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	retry "github.com/giantswarm/retry-go"
	"github.com/wpdirectory/wpdir/internal/client"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/searcher"
	"github.com/wpdirectory/wpdir/internal/ulid"
	"github.com/wpdirectory/wpdir/internal/utils"
)

const (
	archiveURL = "http://downloads.wordpress.org/plugin/%s.latest-stable.zip?nostats=1"
)

// Plugin ...
type Plugin struct {
	Slug                   string             `json:"slug"`
	Name                   string             `json:"name,omitempty"`
	Version                string             `json:"version,omitempty"`
	Author                 string             `json:"author,omitempty"`
	AuthorProfile          string             `json:"author_profile,omitempty"`
	Rating                 int                `json:"rating,omitempty"`
	NumRatings             int                `json:"num_ratings,omitempty"`
	SupportThreads         int                `json:"support_threads,omitempty"`
	SupportThreadsResolved int                `json:"support_threads_resolved,omitempty"`
	ActiveInstalls         int                `json:"active_installs,omitempty"`
	Downloaded             int                `json:"downloaded,omitempty"`
	LastUpdated            string             `json:"last_updated,omitempty"`
	Added                  string             `json:"added,omitempty"`
	Homepage               string             `json:"homepage,omitempty"`
	ShortDescription       string             `json:"short_description,omitempty"`
	DownloadLink           string             `json:"download_link,omitempty"`
	StableTag              string             `json:"stable_tag,omitempty"`
	Status                 status             `json:"status"`
	Searcher               *searcher.Searcher `json:"-"`
	indexed                bool
	sync.RWMutex
}

type status int

const (
	open status = iota
	disabled
	closed
)

// APIResponse ...
type APIResponse struct {
	Name                   string `json:"name"`
	Slug                   string `json:"slug"`
	Version                string `json:"version"`
	Author                 string `json:"author"`
	AuthorProfile          string `json:"author_profile"`
	Rating                 int    `json:"rating"`
	NumRatings             int    `json:"num_ratings"`
	SupportThreads         int    `json:"support_threads"`
	SupportThreadsResolved int    `json:"support_threads_resolved"`
	ActiveInstalls         int    `json:"active_installs"`
	Downloaded             int    `json:"downloaded"`
	LastUpdated            string `json:"last_updated"`
	Added                  string `json:"added"`
	Homepage               string `json:"homepage"`
	ShortDescription       string `json:"short_description"`
	DownloadLink           string `json:"download_link"`
	StableTag              string `json:"stable_tag"`
}

// New returns a new plugin struct.
func New(slug string) *Plugin {

	return &Plugin{
		Slug:   slug,
		Status: closed,
		//Searcher: &searcher.Searcher{},
	}

}

// GetStatus returns the Status as a string
func (p *Plugin) GetStatus() string {
	p.RLock()
	defer p.RUnlock()

	switch p.Status {
	case disabled:
		return "Disabled"
	case closed:
		return "Closed"
	default:
		return "Open"
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
		p.Status = closed
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

	// Main URL Components
	// https://api.wordpress.org/plugins/info/1.1/
	u := &url.URL{
		Scheme: "https",
		Host:   "api.wordpress.org",
		Path:   "plugins/info/1.1/",
	}

	// Query Values
	values := []string{
		"action=plugin_information",
		"request[slug]=" + p.Slug,
		"request[fields][sections]=0",
		"request[fields][description]=0",
		"request[fields][short_description]=1",
		"request[fields][tested]=1",
		"request[fields][requires]=0",
		"request[fields][rating]=1",
		"request[fields][ratings]=0",
		"request[fields][downloaded]=1",
		"request[fields][active_installs]=1",
		"request[fields][last_updated]=1",
		"request[fields][homepage]=1",
		"request[fields][tags]=0",
		"request[fields][donate_link]=0",
		"request[fields][contributors]=0",
		"request[fields][compatibility]=0",
		"request[fields][versions]=0",
		"request[fields][version]=1",
		"request[fields][screenshots]=1",
		"request[fields][stable_tag]=1",
		"request[fields][download_link]=1",
	}

	// Add Query Params to URL and return it as a string
	u.RawQuery = strings.Join(values, "&")
	URL := u.String()

	client := client.GetAPI()

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return data, err
	}

	// Set User-Agent
	req.Header.Set("User-Agent", "wpdirectory/0.1.0")

	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Update ...
func (p *Plugin) Update() error {
	p.Lock()
	defer p.Unlock()

	bytes, err := p.getArchive()
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		p.Status = closed
		return nil
	}

	ref, err := p.processArchive(bytes)
	if err != nil {
		return err
	}

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
func (p *Plugin) processArchive(archive []byte) (*index.IndexRef, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	id := ulid.New()
	dst := filepath.Join(wd, "data", "index", "plugins", id)
	opts := &index.IndexOptions{
		ExcludeDotFiles: true,
	}

	ref, err := index.BuildFromZip(opts, archive, dst, p.Slug)
	if err != nil {
		return nil, err
	}

	return ref, nil
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
