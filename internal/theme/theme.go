package theme

import (
	"encoding/json"
	"errors"
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
	archiveURL = "http://downloads.wordpress.org/theme/%s.latest-stable.zip?nostats=1"
)

// Theme ...
type Theme struct {
	Name          string
	Slug          string
	Version       string
	Author        string
	AuthorProfile string
	Installs      int
	Status        status
	Searcher      *searcher.Searcher
	Stats         *files.Stats `json:"stats,omitempty"`
	indexed       bool
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
	NumRatings             int    `json:"num_ratings"`
	SupportThreads         int    `json:"support_threads"`
	SupportThreadsResolved int    `json:"support_threads_resolved"`
	Downloaded             int    `json:"downloaded"`
	Installs               int    `json:"active_installs"`
	LastUpdated            string `json:"last_updated"`
	Added                  string `json:"added"`
	ShortDescription       string `json:"short_description"`
	DownloadLink           string `json:"download_link"`
}

// New returns a new plugin struct.
func New(slug string) *Theme {

	return &Theme{
		Slug: slug,
	}

}

// GetStatus returns the Status as a string
func (t *Theme) GetStatus() string {
	t.RLock()
	defer t.RUnlock()

	switch t.Status {
	case disabled:
		return "Disabled"
	case closed:
		return "Closed"
	default:
		return "Open"
	}
}

// HasIndex returns the index status
func (t *Theme) HasIndex() bool {
	t.RLock()
	defer t.RUnlock()
	return t.indexed
}

// SetIndexed sets the indexed value
func (t *Theme) SetIndexed(idx bool) {
	t.Lock()
	defer t.Unlock()
	t.indexed = idx
}

// LoadAPIData updates the Plugin struct with data from an HTTP API
func (t *Theme) LoadAPIData() error {
	var data []byte
	var err error

	fetch := func() error {
		data, err = t.getAPIData()
		return err
	}

	err = retry.Do(fetch, retry.Timeout(15*time.Second), retry.MaxTries(3), retry.Sleep(5*time.Second))
	if err != nil || data == nil {
		t.Status = closed
		return err
	}

	err = json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	return nil
}

// getAPIData ...
func (t *Theme) getAPIData() ([]byte, error) {
	var data []byte
	var err error

	api := wporg.NewClient()
	data, err = api.GetInfo("themes", t.Slug)

	return data, err
}

// Update ...
func (t *Theme) Update() error {
	t.Lock()
	defer t.Unlock()

	bytes, err := t.getArchive()
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		t.Status = closed
		return errors.New("No zip file available, Theme is closed")
	}

	ref, stats, err := t.processArchive(bytes)
	if err != nil {
		return err
	}

	// Store File Stats
	t.Stats = stats

	if t.Searcher == nil {
		// New Searcher
		sr, err := searcher.New(ref)
		if err != nil {
			return err
		}
		t.Searcher = sr
	} else {
		// Use Existing Searcher
		idx, err := ref.Open()
		if err != nil {
			return err
		}

		err = t.Searcher.SwapIndexes(idx)
		if err != nil {
			return err
		}
	}

	return nil
}

// getArchive ...
func (t *Theme) getArchive() ([]byte, error) {

	var content []byte
	var err error

	client := client.GetAPI()
	URL := fmt.Sprintf(archiveURL, t.Slug)

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

		log.Printf("Downloading the extension '%s' failed. Response code: %d\n", t.Name, resp.StatusCode)

		return content, err

	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}

	return content, nil

}

// processArchive ...
func (t *Theme) processArchive(archive []byte) (*index.IndexRef, *files.Stats, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	id := ulid.New()
	dst := filepath.Join(wd, "data", "index", "themes", id)
	opts := &index.IndexOptions{
		ExcludeDotFiles: true,
	}

	ref, stats, err := index.BuildFromZip(opts, archive, dst, t.Slug)
	if err != nil {
		return nil, nil, err
	}

	return ref, stats, nil
}

// Save ...
// TODO: Wrap struct to allow locking during Marshal
func (t *Theme) Save() error {
	t.RLock()
	defer t.RUnlock()

	bytes, err := json.Marshal(t)
	if err != nil {
		return err
	}

	db.PutToBucket(t.Slug, bytes, "themes")

	return nil
}
