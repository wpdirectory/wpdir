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
	"github.com/wpdirectory/wpdir/internal/svn"
)

const (
	pluginManagementUser = "plugin-master"
)

// PluginRepo ...
type PluginRepo struct {
	Config *config.Config
	List   map[string]*plugin.Plugin

	Revision    int
	Updated     time.Time
	UpdateQueue chan string
	sync.RWMutex
}

// Len ...
func (pr *PluginRepo) Len() int {
	return len(pr.List)
}

// Rev ...
func (pr *PluginRepo) Rev() int {
	return pr.Revision
}

func (pr *PluginRepo) save() error {
	pr.Lock()
	defer pr.Unlock()

	rev := strconv.Itoa(pr.Revision)

	return db.PutToBucket("plugins", []byte(rev), "repos")
}

func (pr *PluginRepo) load() error {
	pr.RLock()
	defer pr.RUnlock()
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
	pr.Lock()
	defer pr.Unlock()
	_, ok := pr.List[slug]
	return ok
}

// Get ...
func (pr *PluginRepo) Get(slug string) Extension {
	pr.Lock()
	defer pr.Unlock()
	p := pr.List[slug]
	return p
}

// Add ...
func (pr *PluginRepo) Add(slug string) {
	pr.RLock()
	defer pr.RUnlock()
	pr.List[slug] = &plugin.Plugin{
		Slug: slug,
	}
	pr.QueueUpdate(slug)
}

// Set ...
func (pr *PluginRepo) Set(slug string, p *plugin.Plugin) {
	pr.RLock()
	defer pr.RUnlock()
	pr.List[slug] = p
}

// Remove ...
func (pr *PluginRepo) Remove(slug string) {
	pr.RLock()
	defer pr.RUnlock()
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
		pr.List[slug].Status = 1
		return err
	}

	pr.List[slug].SetIndexed(true)
	pr.List[slug].Status = 0

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
			log.Printf("Plugin (%s) Update Failed: %s\n", slug, err)
			//pr.UpdateQueue <- slug
		}
	}
}

// ProcessUpdate ...
func (pr *PluginRepo) ProcessUpdate(slug string) error {
	p := pr.Get(slug).(*plugin.Plugin)
	err := p.LoadAPIData()
	if err != nil {
		p.Status = 1
		return err
	}

	p.Status = 0

	err = p.Update()
	if err != nil {
		p.SetIndexed(false)
		return err
	}
	p.SetIndexed(true)

	p.Save()

	return nil
}

// UpdateList updates our Plugin list.
func (pr *PluginRepo) UpdateList() error {
	// Fetch list from SVN
	// https://plugins.svn.wordpress.org/
	list, err := svn.GetList("plugins", "")
	if err != nil {
		return err
	}

	for _, item := range list {
		if !utf8.Valid([]byte(item.Name)) {
			return errors.New("Plugin slug is not valid utf8")
		}
		if !pr.Exists(item.Name) {
			pr.Add(item.Name)
		}
	}

	return nil
}

// Worker ...
func (pr *PluginRepo) Worker() error {
	updateAPIData := time.NewTicker(time.Hour * 24).C

	checkSVN := time.NewTicker(time.Minute * 5).C

	for {
		select {
		case <-updateAPIData:
			// Update Plugins API Data
			log.Println("Update Pluins API Data.")
		case <-checkSVN:
			// Check SVN for Plugin Updates
			log.Println("Check SVN for Plugin updates.")
		}
	}
}

// LoadExisting ...
func (pr *PluginRepo) LoadExisting() {

	pr.loadDBData()
	pr.loadIndexes()

}

// loadDBData loads all existing Plugin data from the DB.
func (pr *PluginRepo) loadDBData() {
	plugins, err := db.GetAllFromBucket("plugins")
	if err != nil {
		return
	}

	log.Printf("Found %d Plugin(s) in DB.\n", len(plugins))

	for slug, bytes := range plugins {
		var p plugin.Plugin
		err := json.Unmarshal(bytes, &p)
		if err != nil {
			continue
		}
		p.Status = 1
		p.Searcher = &searcher.Searcher{}

		pr.Set(slug, &p)
	}
}

// loadIndexes reads all existing Indexes and attempts to match them to a Plugin.
func (pr *PluginRepo) loadIndexes() {
	indexDir := filepath.Join(pr.Config.WD, "data", "index", "plugins")

	dirs, err := ioutil.ReadDir(indexDir)
	if err != nil {
		log.Printf("Failed to read Plugin index dir: %s\n", err)
		return
	}

	log.Printf("Found %d existing Plugin indexes.\n", len(dirs))

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
	log.Printf("Loaded %d Plugin indexes", loaded)
}

// Summary ...
func (pr *PluginRepo) Summary() *RepoSummary {
	pr.Lock()
	defer pr.Unlock()

	rs := &RepoSummary{
		Revision: pr.Revision,
		Total:    len(pr.List),
		Queue:    len(pr.UpdateQueue),
	}

	for _, p := range pr.List {
		p.Lock()
		if p.Status == 1 {
			rs.Closed++
		}
		p.Unlock()
	}

	return rs
}
