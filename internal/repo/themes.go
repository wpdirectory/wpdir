package repo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/searcher"
	"github.com/wpdirectory/wpdir/internal/theme"
	"github.com/wpdirectory/wporg"
)

const (
	// TODO: Check if this is used like the Plugins Repo.
	themeManagementUser = "theme-master"
)

// ThemeRepo holds data about the Themes SVN Repo.
type ThemeRepo struct {
	Revision    int
	Updated     time.Time
	UpdateQueue chan string

	sync.RWMutex
	List map[string]*theme.Theme

	log *log.Logger
	cfg *config.Config
	api *wporg.Client
}

// Len returns the number of Themes
func (tr *ThemeRepo) Len() uint64 {
	return uint64(len(tr.List))
}

// Rev returns the current Revision
func (tr *ThemeRepo) Rev() int {
	return tr.Revision
}

func (tr *ThemeRepo) save() error {
	tr.RLock()
	defer tr.RUnlock()

	rev := strconv.Itoa(tr.Revision)

	return db.PutToBucket("themes", []byte(rev), "repos")
}

func (tr *ThemeRepo) load() error {
	bytes, err := db.GetFromBucket("themes", "repos")
	if err != nil {
		return err
	}

	rev, err := strconv.Atoi(string(bytes))
	if err != nil || rev != 0 {
		return err
	}
	tr.Revision = rev

	return nil
}

// Exists checks if a Theme exists
func (tr *ThemeRepo) Exists(slug string) bool {
	tr.RLock()
	_, ok := tr.List[slug]
	tr.RUnlock()
	return ok
}

// Get returns a Theme
func (tr *ThemeRepo) Get(slug string) Extension {
	tr.RLock()
	p := tr.List[slug]
	tr.RUnlock()
	return p
}

// Add sets a new Theme
func (tr *ThemeRepo) Add(slug string) {
	tr.Lock()
	tr.List[slug] = theme.New(slug)
	tr.Unlock()
}

// Set ...
func (tr *ThemeRepo) Set(slug string, t *theme.Theme) {
	tr.Lock()
	tr.List[slug] = t
	tr.Unlock()
}

// Remove deletes a current Theme
func (tr *ThemeRepo) Remove(slug string) {
	tr.Lock()
	delete(tr.List, slug)
	tr.Unlock()
}

// UpdateIndex ...
func (tr *ThemeRepo) UpdateIndex(idx *index.Index) error {
	var slug string
	if slug = idx.Ref.Slug; slug == "" {
		// bad index, perhaps delete?
		return errors.New("Index contains empty slug")
	}

	if !tr.Exists(slug) {
		return errors.New("Index does not match an existing theme")
	}

	err := tr.List[slug].Searcher.SwapIndexes(idx)
	if err != nil {
		tr.List[slug].SetIndexed(false)
		tr.List[slug].Status = theme.Closed
		return err
	}

	tr.List[slug].SetIndexed(true)
	tr.List[slug].Status = theme.Open

	return nil
}

// QueueUpdate adds a Theme to the update queue
func (tr *ThemeRepo) QueueUpdate(slug string) {
	tr.UpdateQueue <- slug
}

// UpdateWorker processes updates from the update queue
func (tr *ThemeRepo) UpdateWorker() {
	for {
		slug := <-tr.UpdateQueue
		err := tr.ProcessUpdate(slug)
		if err != nil {
			tr.log.Printf("Theme (%s) Update Failed: %s\n", slug, err)
			//tr.UpdateQueue <- slug
		}
	}
}

// ProcessUpdate ...
func (tr *ThemeRepo) ProcessUpdate(slug string) error {
	if !tr.Exists(slug) {
		tr.Add(slug)
	}
	t := tr.Get(slug).(*theme.Theme)
	err := t.LoadAPIData()
	if err != nil {
		t.Status = theme.Closed
		return err
	}

	err = t.Update()
	if err != nil {
		t.Status = theme.Closed
		t.SetIndexed(false)
		return err
	}

	t.Status = theme.Open
	t.SetIndexed(true)

	t.Save()

	return nil
}

// UpdateList updates our list of themes.
func (tr *ThemeRepo) UpdateList() error {
	// Fetch list from WPOrg API
	list, err := tr.api.GetList("themes")
	if err != nil {
		return err
	}
	tr.log.Printf("Found %d Themes\n", len(list))

	for _, theme := range list {
		if !utf8.Valid([]byte(theme)) {
			return errors.New("Theme slug is not valid UTF8")
		}
		if !tr.Exists(theme) {
			tr.Add(theme)
		}
	}

	return nil
}

// StartWorkers starts up Goroutines to process updates
// Every 15 mins Theme Repos check the changelog for updates
// Every 24 hours all Themes refresh API data
func (tr *ThemeRepo) StartWorkers() {
	// Setup Tickers
	checkChangelog := time.NewTicker(time.Minute * 15).C
	checkAPI := time.NewTicker(time.Hour * 48).C

	go func(tr *ThemeRepo, ticker <-chan time.Time) {
		for {
			select {
			// Check Changlog
			case <-ticker:
				latest, err := tr.api.GetRevision("themes")
				if err != nil {
					tr.log.Printf("Failed getting Themes Repo revision: %s\n", err)
				}
				tr.RLock()
				list, err := tr.api.GetChangeLog("themes", tr.Revision, latest)
				if err != nil {
					tr.log.Printf("Failed getting Themes Changelog: %s\n", err)
					tr.RUnlock()
					continue
				}
				tr.RUnlock()

				for _, slug := range list {
					tr.QueueUpdate(slug)
				}

				tr.Lock()
				tr.Revision = latest
				tr.Unlock()

				err = tr.save()
				if err != nil {
					tr.log.Printf("Failed saving Themes Repo: %s\n", err)
					continue
				}
			}
		}
	}(tr, checkChangelog)

	go func(tr *ThemeRepo, ticker <-chan time.Time) {
		for {
			select {
			// Refresh API Data
			case <-ticker:
				themes, err := tr.api.GetList("themes")
				if err != nil {
					tr.log.Printf("Failed getting Themes list: %s\n", err)
				}
				for _, slug := range themes {
					if !tr.Exists(slug) {
						tr.Add(slug)
					}
					t := tr.Get(slug).(*theme.Theme)
					t.Lock()
					err := t.LoadAPIData()
					if err != nil {
						t.Status = theme.Closed
					}
					t.Unlock()
				}
			}
		}
	}(tr, checkAPI)
}

// LoadExisting ...
func (tr *ThemeRepo) LoadExisting() {
	tr.loadDBData()
	tr.loadIndexes()
}

// loadDBData loads all existing Theme data from the DB.
func (tr *ThemeRepo) loadDBData() {
	themes, err := db.GetAllFromBucket("themes")
	if err != nil {
		return
	}

	tr.log.Printf("Found %d Theme(s) in DB\n", len(themes))

	for slug, bytes := range themes {
		var t theme.Theme
		err := json.Unmarshal(bytes, &t)
		if err != nil {
			continue
		}
		t.Status = theme.Closed
		t.Searcher = &searcher.Searcher{}

		tr.Set(slug, &t)
	}
}

// loadIndexes reads all existing Indexes and attempts to match them to a Theme.
func (tr *ThemeRepo) loadIndexes() {
	indexDir := filepath.Join(tr.cfg.WD, "data", "index", "themes")

	dirs, err := ioutil.ReadDir(indexDir)
	if err != nil {
		tr.log.Printf("Failed to read Theme index dir: %s\n", err)
		return
	}

	tr.log.Printf("Found %d existing Theme indexes\n", len(dirs))

	var loaded int

	for _, dir := range dirs {
		// If not Directory discard.
		if !dir.IsDir() {
			continue
		}

		path := filepath.Join(indexDir, dir.Name())

		// Read Index
		ref, err := index.Read(path)
		if err != nil {
			os.RemoveAll(path)
			continue
		}

		// Create Index
		idx, err := ref.Open()
		if err != nil {
			os.RemoveAll(path)
			continue
		}

		err = tr.UpdateIndex(idx)
		if err != nil {
			os.RemoveAll(path)
			continue
		}
		loaded++
	}
	tr.log.Printf("Loaded %d Theme indexes", loaded)
}

// Summary ...
func (tr *ThemeRepo) Summary() *Summary {
	tr.RLock()
	defer tr.RUnlock()

	rs := &Summary{
		Revision: tr.Revision,
		Total:    len(tr.List),
		Queue:    len(tr.UpdateQueue),
	}

	for _, t := range tr.List {
		t.Lock()
		if t.Status == 1 {
			rs.Closed++
		}
		t.Unlock()
	}

	return rs
}
