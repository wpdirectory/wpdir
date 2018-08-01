package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/wpdirectory/wpdir/internal/client"
	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/filestats"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/ulid"
	"github.com/wpdirectory/wpdir/internal/utils"
	"github.com/wpdirectory/wporg"
)

var (
	archiveURL = "http://downloads.wordpress.org/%s/%s.latest-stable.zip?nostats=1"
)

// Repo holds data about the Plugins SVN Repo.
type Repo struct {
	ExtType     string
	Revision    int
	Updated     time.Time
	UpdateQueue chan UpdateRequest

	sync.RWMutex
	List map[string]*Extension

	log *log.Logger
	cfg *config.Config
	api *wporg.Client
}

// Summary ...
type Summary struct {
	Revision int `json:"revision"`
	Total    int `json:"total"`
	Closed   int `json:"closed"`
	Queue    int `json:"queue"`
}

// New returns a new Repo
func New(c *config.Config, l *log.Logger, t string, rev int) *Repo {
	// Setup HTTP Client
	opt := func(c *wporg.Client) {
		c.HTTPClient = httpClient
	}
	api := wporg.NewClient(opt)

	repo := &Repo{
		cfg:         c,
		log:         l,
		api:         api,
		ExtType:     t,
		Revision:    rev,
		List:        make(map[string]*Extension),
		UpdateQueue: updateQueue,
	}

	// Load Existing Data
	err := repo.load()
	if err != nil {
		l.Printf("Repo (%s) could not load data: %s\n", t, err)
	}

	repo.save()

	return repo
}

// Len ...
func (r *Repo) Len() uint64 {
	return uint64(len(r.List))
}

// Rev ...
func (r *Repo) Rev() int {
	return r.Revision
}

func (r *Repo) save() error {
	r.RLock()
	defer r.RUnlock()

	rev := strconv.Itoa(r.Revision)

	return db.PutToBucket(r.ExtType, []byte(rev), "repos")
}

func (r *Repo) load() error {
	bytes, err := db.GetFromBucket(r.ExtType, "repos")
	if err != nil {
		return err
	}

	rev, err := strconv.Atoi(string(bytes))
	if err != nil || rev == 0 {
		return err
	}

	r.Lock()
	r.Revision = rev
	r.Unlock()
	r.log.Printf("Repo loaded revision: %d\n", rev)

	return nil
}

// Exists ...
func (r *Repo) Exists(slug string) bool {
	r.RLock()
	defer r.RUnlock()
	_, ok := r.List[slug]
	return ok
}

// Get ...
func (r *Repo) Get(slug string) *Extension {
	r.RLock()
	defer r.RUnlock()
	p := r.List[slug]
	return p
}

// Add ...
func (r *Repo) Add(slug string) {
	r.Lock()
	r.List[slug] = NewExt(slug)
	r.Unlock()
}

// Set ...
func (r *Repo) Set(slug string, e *Extension) {
	r.Lock()
	defer r.Unlock()
	r.List[slug] = e
}

// Remove ...
func (r *Repo) Remove(slug string) {
	r.Lock()
	defer r.Unlock()
	delete(r.List, slug)
}

// UpdateIndex ...
func (r *Repo) UpdateIndex(idx *index.Index) error {
	var slug string
	if slug = idx.Ref.Slug; slug == "" {
		// bad index, perhaps delete?
		return errors.New("Index contains empty slug")
	}

	if !r.Exists(slug) {
		return errors.New("Index does not match an existing plugin")
	}

	err := r.List[slug].SwapIndexes(idx)
	if err != nil {
		r.List[slug].SetStatus(Closed)
		return err
	}

	r.List[slug].SetStatus(Open)

	return nil
}

// QueueUpdate ...
func (r *Repo) QueueUpdate(slug string, rev string) {
	revision, err := strconv.Atoi(rev)
	if err != nil {
		r.log.Printf("Revision not a valid int: %s\n", err)
	}
	ur := UpdateRequest{
		Slug:     slug,
		Repo:     r.ExtType,
		Revision: revision,
	}
	r.UpdateQueue <- ur
}

// ProcessUpdate ...
func (r *Repo) ProcessUpdate(slug string, rev int) error {
	if !r.Exists(slug) {
		r.Add(slug)
	}
	e := r.Get(slug)
	err := r.updateMeta(e)
	if err != nil {
		e.SetStatus(Closed)
		return err
	}

	err = r.updateFiles(e)
	if err != nil {
		e.SetStatus(Closed)
		return err
	}

	e.SetStatus(Open)
	r.saveExt(e)

	r.Lock()
	r.Revision = rev
	r.Unlock()

	r.save()

	return nil
}

func (r *Repo) updateMeta(e *Extension) error {
	e.RLock()
	slug := e.Slug
	e.RUnlock()

	bytes, err := r.api.GetInfo(r.ExtType, slug)
	if err != nil {
		return err
	}

	e.Lock()
	defer e.Unlock()
	err = json.Unmarshal(bytes, e)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) updateFiles(e *Extension) error {
	e.RLock()
	slug := e.Slug
	e.RUnlock()

	// Download Extension Archive
	bytes, err := r.getArchive(slug)
	if err != nil {
		return err
	}

	// Index extension using Archive bytes
	ref, files, err := r.generateIndex(bytes, slug)
	if err != nil {
		return err
	}

	// Update File Stats
	e.Lock()
	e.Stats = files
	e.Unlock()

	// Get Index
	idx, err := ref.Open()
	if err != nil {
		return err
	}

	// Update Index
	err = e.SwapIndexes(idx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) getArchive(slug string) ([]byte, error) {
	var content []byte
	var err error

	client := client.GetZip()
	repo := r.ExtType[:len(r.ExtType)-1]
	URL := fmt.Sprintf(archiveURL, repo, slug)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		r.log.Println(err)
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

		log.Printf("Downloading the extension '%s' failed. Response code: %d\n", slug, resp.StatusCode)

		return content, err
	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}

	return content, nil
}

// generateIndex indexes the contents of an archive provided in bytes
func (r *Repo) generateIndex(archive []byte, slug string) (*index.IndexRef, *filestats.Stats, error) {
	id := ulid.New()
	dst := filepath.Join(r.cfg.WD, "data", "index", r.ExtType, id)
	opts := &index.IndexOptions{
		ExcludeDotFiles: true,
	}

	ref, stats, err := index.BuildFromZip(opts, archive, dst, slug)
	if err != nil {
		return nil, nil, err
	}

	return ref, stats, nil
}

// saveExt encodes to JSON and stores in DB
func (r *Repo) saveExt(e *Extension) error {
	e.RLock()
	defer e.RUnlock()

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	db.PutToBucket(e.Slug, bytes, r.ExtType)

	return nil
}

// UpdateList updates our Plugin list.
func (r *Repo) UpdateList(fresh *bool) error {
	// Fetch list from WPOrg API
	list, err := r.api.GetList(r.ExtType)
	if err != nil {
		return err
	}
	r.log.Printf("Found %d %s\n", len(list), r.ExtType)

	revision, err := r.api.GetRevision(r.ExtType)
	if err != nil {
		return err
	}
	rev := strconv.Itoa(revision)

	for _, ext := range list {
		if !utf8.Valid([]byte(ext)) {
			return errors.New("Extension slug is not valid UTF8")
		}
		if !r.Exists(ext) {
			r.Add(ext)
			if *fresh {
				r.QueueUpdate(ext, rev)
			}
		}
	}

	return nil
}

// StartWorkers starts up Goroutines to process updates
// Every 15 mins Plugin Repos check the changelog for updates
// Every 24 hours all Plugins refresh API data
func (r *Repo) StartWorkers() {
	// Setup Tickers
	checkChangelog := time.NewTicker(time.Minute * 15).C
	checkAPI := time.NewTicker(time.Hour * 48).C

	go func(r *Repo, ticker <-chan time.Time) {
		for {
			select {
			// Check Changlog
			case <-ticker:
				latest, err := r.api.GetRevision(r.ExtType)
				if err != nil {
					r.log.Printf("Failed getting %s Repo revision: %s\n", r.ExtType, err)
				}
				r.RLock()
				list, err := r.api.GetChangeLog(r.ExtType, r.Revision, latest)
				if err != nil {
					r.log.Printf("Failed getting %s Changelog: %s\n", r.ExtType, err)
					r.RUnlock()
					continue
				}
				r.RUnlock()

				for _, ext := range list {
					r.QueueUpdate(string(ext[0]), string(ext[1]))
				}
			}
		}
	}(r, checkChangelog)

	go func(r *Repo, ticker <-chan time.Time) {
		for {
			select {
			// Refresh API Data
			case <-ticker:
				exts, err := r.api.GetList(r.ExtType)
				if err != nil {
					r.log.Printf("Failed getting %s list: %s\n", r.ExtType, err)
				}
				for _, slug := range exts {
					if !r.Exists(slug) {
						r.Add(slug)
					}
					e := r.Get(slug)
					err := r.updateMeta(e)
					if err != nil {
						e.SetStatus(Closed)
					}
				}
			}
		}
	}(r, checkAPI)
}

// LoadExisting loading data from DB and then Indexes
func (r *Repo) LoadExisting() {
	r.loadDBData()
	r.loadIndexes()
}

// loadDBData loads all existing Plugin data from the DB
func (r *Repo) loadDBData() {
	exts, err := db.GetAllFromBucket(r.ExtType)
	if err != nil {
		return
	}

	r.log.Printf("Found %d %s in DB\n", len(exts), r.ExtType)

	for slug, bytes := range exts {
		var e Extension
		err := json.Unmarshal(bytes, &e)
		if err != nil {
			continue
		}
		e.Status = Closed

		r.Set(slug, &e)
	}
}

// loadIndexes reads all existing Indexes and attempts to match them to a Plugin.
func (r *Repo) loadIndexes() {
	indexDir := filepath.Join(r.cfg.WD, "data", "index", r.ExtType)

	dirs, err := ioutil.ReadDir(indexDir)
	if err != nil {
		r.log.Printf("Failed to read %s index dir: %s\n", r.ExtType, err)
		return
	}

	r.log.Printf("Found %d existing %s indexes\n", len(dirs), r.ExtType)

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

		err = r.UpdateIndex(idx)
		if err != nil {
			os.RemoveAll(path)
			continue
		}
		loaded++
	}
	r.log.Printf("Loaded %d/%d indexes", loaded, len(dirs))
}

// Summary ...
func (r *Repo) Summary() *Summary {
	r.RLock()
	defer r.RUnlock()

	rs := &Summary{
		Revision: r.Revision,
		Total:    len(r.List),
		Queue:    len(r.UpdateQueue),
	}

	for _, e := range r.List {
		e.RLock()
		if e.Status == Closed {
			rs.Closed++
		}
		e.RUnlock()
	}

	return rs
}
