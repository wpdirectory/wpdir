package theme

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	retry "github.com/giantswarm/retry-go"
	"github.com/wpdirectory/wpdir/internal/client"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/searcher"
	"github.com/wpdirectory/wpdir/internal/utils"
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
	t.Lock()
	defer t.Unlock()
	return t.indexed
}

// SetIndexed sets the indexed value
func (t *Theme) SetIndexed(idx bool) {
	t.RLock()
	defer t.RUnlock()
	t.indexed = idx
}

// LoadAPIData updates the Plugin struct with data from an HTTP API
func (t *Theme) LoadAPIData() error {

	var resp *APIResponse
	var err error

	fetch := func() error {
		resp, err = t.getAPIData()
		return err
	}

	if resp == nil {
		t.Status = closed
		return nil
	}

	err = retry.Do(fetch, retry.Timeout(15*time.Second), retry.MaxTries(3), retry.Sleep(5*time.Second))
	if err != nil {
		t.Status = disabled
		return err
	}

	// Update from API data
	t.Name = resp.Name
	t.Version = resp.Version
	t.Author = resp.Author
	t.AuthorProfile = resp.AuthorProfile
	t.Installs = resp.Installs
	t.Status = open

	return nil

}

// getAPIData ...
func (t *Theme) getAPIData() (*APIResponse, error) {

	var result *APIResponse

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
		"request[slug]=" + t.Slug,
		"request[fields][sections]=0",
		"request[fields][description]=0",
		"request[fields][short_description]=1",
		"request[fields][tested]=1",
		"request[fields][requires]=1",
		"request[fields][rating]=1",
		"request[fields][ratings]=1",
		"request[fields][downloaded]=1",
		"request[fields][active_installs]=1",
		"request[fields][last_updated]=1",
		"request[fields][homepage]=1",
		"request[fields][tags]=1",
		"request[fields][donate_link]=0",
		"request[fields][contributors]=0",
		"request[fields][compatibility]=1",
		"request[fields][versions]=0",
		"request[fields][version]=1",
		"request[fields][screenshots]=1",
		"request[fields][stable_tag]=1",
		"request[fields][download_link]=1",
	}

	// Add Query Params to URL and return it as a string
	u.RawQuery = strings.Join(values, "&")
	URL := u.String()

	// Make the HTTP request
	response, err := http.Get(URL)
	if err != nil {
		return result, err
	}

	defer response.Body.Close()
	bodyByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	json.Unmarshal([]byte(bodyByte), &result)

	return result, nil

}

// Update ...
func (t *Theme) Update() error {
	t.RLock()
	defer t.RUnlock()

	bytes, err := t.getArchive()
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		t.Status = closed
		return errors.New("No zip file available, Theme is closed")
	}

	ref, err := t.processArchive(bytes)
	if err != nil {
		return err
	}

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
func (t *Theme) processArchive(archive []byte) (*index.IndexRef, error) {
	dst := filepath.Join()
	opts := &index.IndexOptions{
		ExcludeDotFiles: true,
	}

	ref, err := index.BuildFromZip(opts, archive, dst, t.Slug)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

// Save ...
func (t *Theme) Save() error {
	t.Lock()
	defer t.Unlock()

	bytes, err := json.Marshal(t)
	if err != nil {
		return err
	}

	db.PutToBucket(t.Slug, bytes, "themes")

	return nil
}
