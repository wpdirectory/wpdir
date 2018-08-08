package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sort"
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
	"github.com/wcharczuk/go-chart"
) 

var (
	archiveURL = "http://downloads.wordpress.org/%s/%s.latest-stable.zip?nostats=1"
)

// Repo holds data about the Plugins SVN Repo.
type Repo struct {
	ExtType     string             `json:"type"`
	Revision    int                `json:"revision"`
	Updated     time.Time          `json:"updated"`
	Total       int                `json:"total"`
	Closed      int                `json:"closed"`
	UpdateQueue chan UpdateRequest `json:"-"`

	List map[string]*Extension `json:"-"`
	sync.RWMutex

	log *log.Logger
	cfg *config.Config
	api *wporg.Client
}

// New returns a new Repo
func New(c *config.Config, l *log.Logger, t string, rev int) *Repo {
	// Setup HTTP Client
	opt := func(c *wporg.Client) {
		c.HTTPClient = client.GetAPI()
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

// Len returns the number of extensions in the Repository
func (r *Repo) Len() uint64 {
	r.RLock()
	defer r.RUnlock()

	return uint64(len(r.List))
}

// GetRev returns the Revision of the the Repository
func (r *Repo) GetRev() int {
	r.RLock()
	defer r.RUnlock()

	return r.Revision
}

// SetRev sets the Revision of the the Repository
func (r *Repo) SetRev(new int) {
	r.Lock()
	// Check this is a progression
	if new > r.Revision {
		r.Revision = new
		r.Updated = time.Now()
	}
	r.Unlock()
}

// save stores the Repo data in the DB
func (r *Repo) save() error {
	r.RLock()
	defer r.RUnlock()

	rev := strconv.Itoa(r.Revision)

	return db.PutToBucket(r.ExtType, []byte(rev), "repos")
}

// load gets the Repo data from DB
func (r *Repo) load() error {
	b, err := db.GetFromBucket(r.ExtType, "repos")
	if err != nil {
		return err
	}

	rev, err := strconv.Atoi(string(b))
	if err != nil || rev == 0 {
		return err
	}

	r.SetRev(rev)
	r.log.Printf("Repo loaded revision: %d\n", rev)

	return nil
}

// Exists checks if an extension exists in the Repo
func (r *Repo) Exists(slug string) bool {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.List[slug]
	return ok
}

// Get returns a pointer to an Extension
func (r *Repo) Get(slug string) *Extension {
	r.RLock()
	defer r.RUnlock()
	p := r.List[slug]
	return p
}

// Add creates a new Extension in the Repo
func (r *Repo) Add(slug string) {
	r.Lock()
	r.List[slug] = newExt(slug)
	r.Total++
	r.Closed++
	r.Unlock()
}

// Set loads the provided Extension into the Repo
func (r *Repo) Set(slug string, e *Extension) {
	r.Lock()
	defer r.Unlock()

	r.List[slug] = e
}

// Remove deletes an Extension from the Repo
func (r *Repo) Remove(slug string) {
	r.Lock()
	defer r.Unlock()

	r.Total--
	delete(r.List, slug)
}

// SetStatus sets the Extension Status
func (r *Repo) SetStatus(e *Extension, s status) {
	e.Lock()
	defer e.Unlock()

	r.Lock()
	if e.Status == Closed {
		r.Closed--
	}
	if s == Closed {
		r.Closed++
	}
	r.Unlock()

	e.Status = s
}

// UpdateIndex updates the index held by an Extension
func (r *Repo) UpdateIndex(idx *index.Index) error {
	var slug string
	if slug = idx.Ref.Slug; slug == "" {
		// bad index, perhaps delete?
		return errors.New("Index contains empty slug")
	}

	if !r.Exists(slug) {
		return errors.New("Index does not match an existing plugin")
	}

	// Swap the old index for the new
	err := r.List[slug].SwapIndexes(idx)
	if err != nil {
		r.SetStatus(r.Get(slug), Closed)
		return err
	}

	r.SetStatus(r.Get(slug), Open)

	return nil
}

// QueueUpdate adds a request to the Update Queue
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

// ProcessUpdate performs an update
// Updates Meta data and files
func (r *Repo) ProcessUpdate(slug string, rev int) error {
	if !r.Exists(slug) {
		r.Add(slug)
	}
	e := r.Get(slug)

	// Get latest API info
	err := r.updateMeta(e)
	if err != nil {
		r.SetStatus(e, Closed)
		return err
	}

	// Get latest files
	err = r.updateFiles(e)
	if err != nil {
		r.SetStatus(e, Closed)
		return err
	}

	r.SetStatus(e, Open)
	r.saveExt(e)

	r.SetRev(rev)
	r.save()

	//err = utils.RemoveContents(filepath.Join(r.cfg.WD, "tmp"))
	//if err != nil {
		//r.log.Printf("Failed to clean tmp dir: %s\n", err)
	//}

	return nil
}

// updateMeta updates the Info held for the Extension
func (r *Repo) updateMeta(e *Extension) error {
	e.RLock()
	slug := e.Slug
	e.RUnlock()

	// Fetch API Response
	b, err := r.api.GetInfo(r.ExtType, slug)
	if err != nil {
		return err
	}

	e.Lock()
	defer e.Unlock()
	err = json.Unmarshal(b, e)
	if err != nil {
		return err
	}

	return nil
}

// updateFiles updates the files and index for the Extension
func (r *Repo) updateFiles(e *Extension) error {
	e.RLock()
	slug := e.Slug
	e.RUnlock()

	// Download Extension Archive
	b, err := r.getArchive(slug)
	if err != nil {
		return err
	}

	// Index extension using Archive bytes
	ref, files, err := r.generateIndex(b, slug)
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

// getArchive fetches the latest archive containing Extension files
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
	agent := r.cfg.Name + "/" + r.cfg.Version
	req.Header.Set("User-Agent", agent)

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

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	db.PutToBucket(e.Slug, b, r.ExtType)

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

	// Get latest Revision
	revision, err := r.api.GetRevision(r.ExtType)
	if err != nil {
		return err
	}
	rev := strconv.Itoa(revision)

	for _, ext := range list {
		if !utf8.Valid([]byte(ext)) {
			r.log.Printf("Extension slug is not valid UTF8: %s\n", ext)
			continue
		}
		if !r.Exists(ext) {
			r.Add(ext)
		}
		// If fresh start we should update all Extensions
		if *fresh || r.Revision == 0 {
			r.QueueUpdate(ext, rev)
		}
	}

	return nil
}

// StartWorkers starts up Goroutines to process updates
// Every 15 mins Plugin Repos check the changelog for updates
// Every 24 hours all Plugins refresh API data
// TODO: All a job to clean out files created in temp dir
func (r *Repo) StartWorkers() {
	// Setup Tickers
	checkChangelog := time.NewTicker(time.Minute * 15).C
	checkAPI := time.NewTicker(time.Hour * 48).C

	// Fetch the Changelog to get a list of Extensions to update
	go func(r *Repo, ticker <-chan time.Time) {
		for {
			select {
			// Check Changlog
			case <-ticker:
				// Skip if the update queue is not empty
				if len(r.UpdateQueue) > 0 {
					continue
				}

				latest, err := r.api.GetRevision(r.ExtType)
				if err != nil {
					r.log.Printf("Failed getting %s Repo revision: %s\n", r.ExtType, err)
				}
				r.RLock()
				log, err := r.api.GetChangeLog(r.ExtType, r.Revision, latest)
				if err != nil {
					r.log.Printf("Failed getting %s Changelog: %s\n", r.ExtType, err)
					r.RUnlock()
					continue
				}

				// If no changes skip
				if len(log) == 0 {
					r.log.Printf("No new %s updates since: %d\n", r.ExtType, r.Revision)
					r.RUnlock()
					continue
				}
				r.RUnlock()

				// Remove Duplicates
				// Save most recent revision
				list := make(map[string]string)
				for _, update := range log {
					list[update[0]] = update[1]
				}

				// Queue Updates
				for k, v := range list {
					r.QueueUpdate(k, v)
				}

				r.log.Printf("%d %s added to the update queue\n", len(list), r.ExtType)
			}
		}
	}(r, checkChangelog)

	// Refresh Extension API data
	// It will add missing Extensions to the Repo, but these should be caught in the Changelog above
	go func(r *Repo, ticker <-chan time.Time) {
		for {
			select {
			// Refresh API Data
			case <-ticker:
				exts, err := r.api.GetList(r.ExtType)
				if err != nil {
					r.log.Printf("Failed getting %s list: %s\n", r.ExtType, err)
				}
				for _, ext := range exts {
					if !r.Exists(ext) {
						r.Add(ext)
					}
					e := r.Get(ext)
					err := r.updateMeta(e)
					if err != nil {
						r.SetStatus(e, Closed)
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

	r.Total = 0
	r.Closed = 0

	for _, ext := range r.List {
		r.Total++
		if ext.Status == Closed {
			r.Closed++
		}
	}
}

// loadDBData loads all existing Plugin data from the DB
func (r *Repo) loadDBData() {
	exts, err := db.GetAllFromBucket(r.ExtType)
	if err != nil {
		return
	}

	r.log.Printf("Found %d %s in DB\n", len(exts), r.ExtType)

	for slug, b := range exts {
		var e Extension
		err := json.Unmarshal(b, &e)
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

		r.SetStatus(r.Get(ref.Slug), Open)

		loaded++
	}
	r.log.Printf("Loaded %d/%d indexes", loaded, len(dirs))
}

// GetInstallsChart ...
func (r *Repo) GetInstallsChart() string {
	yr, m, _ := time.Now().Date()
	key := fmt.Sprintf("installs_%s_%d_%d", r.ExtType, yr, m)

	b, err := db.GetFromBucket(key, "charts")
	if err != nil {
		b = r.GenerateInstallsChart()
	}

	return string(b)
}

// generateInstallsChart ...
func (r *Repo) GenerateInstallsChart() []byte {
	r.RLock()
	list := r.List
	r.RUnlock()
	var x, y []float64
	i := 1
	for _, ext := range list {
		x = append(x, 1.00 * float64(i))
		ext.RLock()
		y = append(y, float64(ext.ActiveInstalls))
		ext.RUnlock()
		i++
	}

	sort.Float64s(y)

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      r.ExtType,
			NameStyle: chart.StyleShow(),
			Style:     chart.Style{
				Show: false,
			},
		},
		YAxis: chart.YAxis{
			Name:      "Installs",
			NameStyle: chart.Style{
				Show: true,
				FontSize: 16,
			},
			Style:     chart.Style{
				Show: true,
				FontSize: 16,
			},
			ValueFormatter: func(v interface{}) string {
				if v, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%d", int(v))
				}
				return ""
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: y,
			},
		},
	}

	b := bytes.NewBuffer([]byte{})
	graph.Render(chart.SVG, b)

	// Alter SVG Tag
	str := string(b.Bytes())
	str = strings.Replace(str, "width=\"1024\"", "", 1)
	str = strings.Replace(str, "height=\"400\"", "viewBox=\"0 0 1024 400\"", 1)

	yr, m, _ := time.Now().Date()
	key := fmt.Sprintf("installs_%s_%d_%d", r.ExtType, yr, m)
	err := db.PutToBucket(key, []byte(str), "charts")
	if err != nil {
		r.log.Printf("Error saving %s installs Chart: %s\n", r.ExtType, err)
	}

	return []byte(str)
}

// GetSizeChart ...
func (r *Repo) GetSizeChart() string {
	yr, m, _ := time.Now().Date()
	key := fmt.Sprintf("size_%s_%d_%d", r.ExtType, yr, m)

	b, err := db.GetFromBucket(key, "charts")
	if err != nil {
		b = r.GenerateSizeChart()
	}

	return string(b)
}

// generateSizeChart ...
func (r *Repo) GenerateSizeChart() []byte {
	r.RLock()
	list := r.List
	r.RUnlock()
	var x, y1, y2 []float64
	i := 1
	for _, ext := range list {
		if ext.Stats == nil {
			i++
			continue
		}
		x = append(x, 1.00 * float64(i))
		ext.RLock()
		size := (float64(ext.Stats.TotalSize) / 1024) / 1024
		y1 = append(y1, size)
		y2 = append(y2, float64(ext.Stats.TotalFiles))
		ext.RUnlock()
		i++
	}

	sort.Float64s(y1)
	sort.Float64s(y2)

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      r.ExtType,
			NameStyle: chart.StyleShow(),
			Style:     chart.Style{
				Show: false,
			},
		},
		YAxis: chart.YAxis{
			Name:      "Size (MB)",
			NameStyle: chart.Style{
				Show: true,
				FontSize: 16,
			},
			Style:     chart.Style{
				Show: true,
				FontSize: 16,
			},
			ValueFormatter: func(v interface{}) string {
				if v, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%d", int(v))
				}
				return ""
			},
		},
		YAxisSecondary: chart.YAxis{
			Name:      "Files",
			NameStyle: chart.Style{
				Show: true,
				FontSize: 16,
			},
			Style:     chart.Style{
				Show: true,
				FontSize: 16,
			},
			ValueFormatter: func(v interface{}) string {
				if v, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%.0f", v)
				}
				return ""
			},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 40,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Total Size (MB)",
				XValues: x,
				YValues: y1,
			},
			chart.ContinuousSeries{
				Name:    "Total Files",
				YAxis:   chart.YAxisSecondary,
				XValues: x,
				YValues: y2,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(
			&graph,
			chart.Style{
				Show: true,
				FontSize: 16,
			},
		),
	}

	b := bytes.NewBuffer([]byte{})
	graph.Render(chart.SVG, b)

	// Alter SVG Tag
	str := string(b.Bytes())
	str = strings.Replace(str, "width=\"1024\"", "", 1)
	str = strings.Replace(str, "height=\"400\"", "viewBox=\"0 0 1024 400\"", 1)

	yr, m, _ := time.Now().Date()
	key := fmt.Sprintf("size_%s_%d_%d", r.ExtType, yr, m)
	err := db.PutToBucket(key, []byte(str), "charts")
	if err != nil {
		r.log.Printf("Error saving %s size Chart: %s\n", r.ExtType, err)
	}

	return []byte(str)
}