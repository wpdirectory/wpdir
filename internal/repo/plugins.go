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
	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/searcher"
	"github.com/wpdirectory/wporg"
)

const (
	// This is the user which administrates the Plugins Repo
	// Adds the initial folder structure for approved plugins
	// TODO: We might be able to use this to avoid needless updates
	pluginManagementUser = "plugin-master"
)

// PluginRepo holds data about the Plugins SVN Repo.
type PluginRepo struct {
	Revision    int
	Updated     time.Time
	UpdateQueue chan string

	sync.RWMutex
	List map[string]*plugin.Plugin

	log *log.Logger
	cfg *config.Config
	api *wporg.Client
}

// Len ...
func (pr *PluginRepo) Len() uint64 {
	return uint64(len(pr.List))
}

// Rev ...
func (pr *PluginRepo) Rev() int {
	return pr.Revision
}

func (pr *PluginRepo) save() error {
	pr.RLock()
	defer pr.RUnlock()

	rev := strconv.Itoa(pr.Revision)

	return db.PutToBucket("plugins", []byte(rev), "repos")
}

func (pr *PluginRepo) load() error {
	//pr.Lock()
	//defer pr.Unlock()

	bytes, err := db.GetFromBucket("plugins", "repos")
	if err != nil {
		return err
	}

	rev, err := strconv.Atoi(string(bytes))
	if err != nil || rev != 0 {
		return err
	}
	pr.Revision = rev

	return nil
}

// Exists ...
func (pr *PluginRepo) Exists(slug string) bool {
	pr.RLock()
	defer pr.RUnlock()
	_, ok := pr.List[slug]
	return ok
}

// Get ...
func (pr *PluginRepo) Get(slug string) Extension {
	pr.RLock()
	defer pr.RUnlock()
	p := pr.List[slug]
	return p
}

// Add ...
func (pr *PluginRepo) Add(slug string) {
	pr.Lock()
	pr.List[slug] = &plugin.Plugin{
		Slug:   slug,
		Status: plugin.Closed,
	}
	pr.Unlock()
}

// Set ...
func (pr *PluginRepo) Set(slug string, p *plugin.Plugin) {
	pr.Lock()
	defer pr.Unlock()
	pr.List[slug] = p
}

// Remove ...
func (pr *PluginRepo) Remove(slug string) {
	pr.Lock()
	defer pr.Unlock()
	delete(pr.List, slug)
}

// UpdateIndex ...
func (pr *PluginRepo) UpdateIndex(idx *index.Index) error {
	var slug string
	if slug = idx.Ref.Slug; slug == "" {
		// bad index, perhaps delete?
		return errors.New("Index contains empty slug")
	}

	if !pr.Exists(slug) {
		return errors.New("Index does not match an existing plugin")
	}

	err := pr.List[slug].Searcher.SwapIndexes(idx)
	if err != nil {
		pr.List[slug].SetIndexed(false)
		pr.List[slug].Status = plugin.Closed
		return err
	}

	pr.List[slug].SetIndexed(true)
	pr.List[slug].Status = plugin.Open

	return nil
}

// QueueUpdate ...
func (pr *PluginRepo) QueueUpdate(slug string) {
	pr.UpdateQueue <- slug
}

// UpdateWorker ...
func (pr *PluginRepo) UpdateWorker() {
	for {
		slug := <-pr.UpdateQueue
		err := pr.ProcessUpdate(slug)
		if err != nil {
			pr.log.Printf("Plugin (%s) Update Failed: %s\n", slug, err)
			//pr.UpdateQueue <- slug
		}
	}
}

// ProcessUpdate ...
func (pr *PluginRepo) ProcessUpdate(slug string) error {
	p := pr.Get(slug).(*plugin.Plugin)
	err := p.LoadAPIData()
	if err != nil {
		p.Status = plugin.Closed

		return err
	}

	err = p.Update()
	if err != nil {
		p.Status = plugin.Closed
		p.SetIndexed(false)
		return err
	}
	p.Status = plugin.Open
	p.SetIndexed(true)

	p.Save()

	return nil
}

// UpdateList updates our Plugin list.
func (pr *PluginRepo) UpdateList() error {
	// Fetch list from WPOrg API
	list, err := pr.api.GetList("plugins")
	if err != nil {
		return err
	}
	pr.log.Printf("Found %d Plugins\n", len(list))

	for _, plugin := range list {
		if !utf8.Valid([]byte(plugin)) {
			return errors.New("Plugin slug is not valid UTF8")
		}
		if !pr.Exists(plugin) {
			pr.Add(plugin)
		}
	}

	return nil
}

// StartWorkers starts up Goroutines to process updates
// Every 15 mins Plugin Repos check the changelog for updates
// Every 24 hours all Plugins refresh API data
func (pr *PluginRepo) StartWorkers() {
	// Setup Tickers
	checkChangelog := time.NewTicker(time.Minute * 15).C
	checkAPI := time.NewTicker(time.Hour * 24).C

	go func(ticker <-chan time.Time) {
		for {
			select {
			// Check Changlog
			case <-ticker:
				latest, err := pr.api.GetRevision("plugins")
				if err != nil {
					pr.log.Printf("Failed getting Plugins Repo revision: %s\n", err)
				}
				pr.RLock()
				defer pr.RUnlock()
				list, err := pr.api.GetChangeLog("plugins", pr.Revision, latest)
				if err != nil {
					pr.log.Printf("Failed getting Plugins Changelog: %s\n", err)
				}
				for _, slug := range list {
					pr.QueueUpdate(slug)
				}
				err = pr.save()
				pr.log.Printf("Failed saving Plugins Repo: %s\n", err)
			}
		}
	}(checkChangelog)

	go func(ticker <-chan time.Time) {
		for {
			select {
			// Refresh API Data
			case <-ticker:
				plugins, err := pr.api.GetList("plugins")
				if err != nil {
					pr.log.Printf("Failed getting Plugins list: %s\n", err)
				}
				for _, slug := range plugins {
					p := pr.Get(slug).(*plugin.Plugin)
					p.RLock()
					defer p.RUnlock()
					err := p.LoadAPIData()
					if err != nil {
						p.Status = plugin.Closed
					}
					p.Status = plugin.Closed
				}
			}
		}
	}(checkAPI)
}

// LoadExisting loading data from DB and then Indexes
func (pr *PluginRepo) LoadExisting() {
	pr.loadDBData()
	pr.loadIndexes()
}

// loadDBData loads all existing Plugin data from the DB
func (pr *PluginRepo) loadDBData() {
	plugins, err := db.GetAllFromBucket("plugins")
	if err != nil {
		return
	}

	pr.log.Printf("Found %d Plugin(s) in DB\n", len(plugins))

	for slug, bytes := range plugins {
		var p plugin.Plugin
		err := json.Unmarshal(bytes, &p)
		if err != nil {
			continue
		}
		p.Status = plugin.Closed
		p.Searcher = &searcher.Searcher{}

		pr.Set(slug, &p)
	}
}

// loadIndexes reads all existing Indexes and attempts to match them to a Plugin.
func (pr *PluginRepo) loadIndexes() {
	indexDir := filepath.Join(pr.cfg.WD, "data", "index", "plugins")

	dirs, err := ioutil.ReadDir(indexDir)
	if err != nil {
		pr.log.Printf("Failed to read Plugin index dir: %s\n", err)
		return
	}

	pr.log.Printf("Found %d existing Plugin indexes\n", len(dirs))

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

		err = pr.UpdateIndex(idx)
		if err != nil {
			os.RemoveAll(path)
			continue
		}
		loaded++
	}
	pr.log.Printf("Loaded %d Plugin indexes", loaded)
}

// Summary ...
func (pr *PluginRepo) Summary() *Summary {
	pr.RLock()
	defer pr.RUnlock()

	rs := &Summary{
		Revision: pr.Revision,
		Total:    len(pr.List),
		Queue:    len(pr.UpdateQueue),
	}

	for _, p := range pr.List {
		p.RLock()
		if p.Status == 1 {
			rs.Closed++
		}
		p.RUnlock()
	}

	return rs
}
