package search

import (
	"errors"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/metrics"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/search/queue"
	"github.com/wpdirectory/wpdir/internal/ulid"
)

// Manager controls the processing and storage of searches
type Manager struct {
	Queue   *queue.Queue
	List    map[string]*Search
	Plugins *repo.Repo
	Themes  *repo.Repo
	limit   int
	sync.RWMutex
}

// NewManager returns a new SearchManager struct
func NewManager(limit int) *Manager {
	return &Manager{
		Queue: queue.New(100),
		List:  make(map[string]*Search),
		limit: limit,
	}
}

// Get ...
func (sm *Manager) Get(ID string) Search {
	sm.RLock()
	defer sm.RUnlock()
	s := sm.List[ID]
	return *s
}

// Set ...
func (sm *Manager) Set(s *Search) {
	sm.Lock()
	defer sm.Unlock()
	_, ok := sm.List[s.ID]
	if !ok {
		sm.List[s.ID] = s
	}
}

// Exists ...
func (sm *Manager) Exists(ID string) bool {
	sm.RLock()
	defer sm.RUnlock()
	_, ok := sm.List[ID]
	return ok
}

// Empty ...
func (sm *Manager) Empty() error {
	return db.DeleteSearches()
}

// NewSearch ...
func (sm *Manager) NewSearch(sr Request) string {
	sm.Lock()
	defer sm.Unlock()

	ID := ulid.New()
	sm.List[ID] = &Search{
		ID:      ID,
		Input:   sr.Input,
		Repo:    sr.Repo,
		Private: sr.Private,
		Options: &sr.Opts,
		Status:  Queued,
	}

	sm.Queue.Add(ID)

	return ID
}

// Worker ...
func (sm *Manager) Worker() {

	for {
		searchID := sm.Queue.Get()
		err := sm.processSearch(searchID)
		if err != nil {
			log.Printf("Searched failed: %s\n", err)
		}
	}

}

// SummaryList ...
type SummaryList struct {
	List  map[string]*Result
	Total uint64
	sync.RWMutex
}

// MatchList ...
type MatchList struct {
	List map[string]*Matches
	sync.RWMutex
}

// processSearch ...
func (sm *Manager) processSearch(ID string) error {
	start := time.Now()

	sm.RLock()
	srch, ok := sm.List[ID]
	sm.RUnlock()
	if !ok {
		return errors.New("No search found")
	}

	opts := &index.SearchOptions{
		Offset:         0,
		Limit:          0,
		LinesOfContext: 2,
		IgnoreCase:     false,
		IgnoreComments: false,
	}

	var total, current, totalMatches uint64
	var input string

	sm.RLock()
	input = srch.Input
	searchID := srch.ID
	sm.RUnlock()

	sm.Lock()
	srch.Started = time.Now().Format(time.RFC3339)
	srch.Status = Started
	srch.Matches = 0
	srch.Options = &Options{
		IgnoreCase:     opts.IgnoreCase,
		LinesOfContext: uint32(opts.LinesOfContext),
		IgnoreComments: opts.IgnoreComments,
		Offset:         uint32(opts.Offset),
		Limit:          uint32(opts.Limit),
	}
	sm.Unlock()

	sum := &SummaryList{
		List:  make(map[string]*Result),
		Total: 0,
	}
	matchList := &MatchList{
		List: make(map[string]*Matches),
	}

	limiter := make(chan struct{}, sm.limit)

	var wg sync.WaitGroup

	var r *repo.Repo

	switch srch.Repo {
	case "plugins":
		r = sm.Plugins
		break
	case "themes":
		r = sm.Themes
		break
	default:
		return errors.New("Not a valid repository name")
	}

	r.RLock()
	total = uint64(len(r.List))
	list := r.List
	srch.Revision = uint32(r.Revision)
	r.RUnlock()
	for _, e := range list {
		// Limit to 100000 matches
		if totalMatches > 100000 {
			break
		}
		current++
		sm.Lock()
		srch.Progress = uint32(math.Round((float64(current) / float64(total)) * 100.00))
		srch.Matches = uint32(totalMatches)
		sm.Unlock()
		if e.Status != repo.Open {
			continue
		}
		wg.Add(1)
		limiter <- struct{}{}

		go func(e *repo.Extension, input string, sum *SummaryList, matchlist *MatchList, totalMatches *uint64, wg *sync.WaitGroup) {
			e.RLock()
			defer e.RUnlock()
			resp, err := e.Search(input, e.Slug, opts)
			if err != nil || len(resp.Matches) == 0 {
				wg.Done()
				<-limiter
				return
			}
			var eMatches uint64
			for i := 0; i < len(resp.Matches); i++ {
				if resp.Matches[i].Matches == nil {
					continue
				}
				eMatches = uint64(len(resp.Matches[i].Matches))
				atomic.AddUint64(totalMatches, eMatches)
				ms := &Matches{}
				for j := 0; j < len(resp.Matches[i].Matches); j++ {
					text := resp.Matches[i].Matches[j].Line
					if len(text) > 100 {
						text = text[0:100]
					}
					m := &Match{
						Slug:     e.Slug,
						File:     resp.Matches[i].Filename,
						LineNum:  uint32(resp.Matches[i].Matches[j].LineNumber),
						LineText: text,
					}
					ms.List = append(ms.List, m)
				}
				matchList.Lock()
				matchList.List[e.Slug] = ms
				matchList.Unlock()
			}
			r := &Result{
				Slug:           e.Slug,
				Name:           e.Name,
				Version:        e.Version,
				Homepage:       e.Homepage,
				ActiveInstalls: uint32(e.ActiveInstalls),
				Matches:        uint32(eMatches),
			}
			sum.Lock()
			sum.List[e.Slug] = r
			sum.Unlock()
			wg.Done()
			<-limiter
		}(e, input, sum, matchList, &totalMatches, &wg)
	}

	wg.Wait()

	summary := &Summary{
		List:  make(map[string]*Result),
		Total: sum.Total,
	}
	for key, result := range sum.List {
		summary.List[key] = result
	}

	bytes, err := summary.Marshal()
	if err != nil {
		return errors.New("Failed Marshalling Summary")
	}
	err = db.SaveSummary(searchID, bytes)
	if err != nil {
		return errors.New("Failed Saving Summary to DB")
	}

	// Create new map with Marshal bytes so that the DB can be run as a transaction
	mlist := make(map[string][]byte, len(matchList.List))
	for slug, matches := range matchList.List {
		bytes, err = matches.Marshal()
		if err != nil {
			return errors.New("Failed Marshalling Matches")
		}
		mlist[slug] = bytes
	}
	err = db.SaveMatches(searchID, mlist)
	if err != nil {
		return errors.New("Failed Saving Summary to DB")
	}

	sm.Lock()
	srch.Completed = time.Now().Format(time.RFC3339)
	srch.Status = Completed
	srch.Matches = uint32(totalMatches)
	sm.Unlock()

	sm.RLock()
	s := *srch
	bytes, err = srch.Marshal()
	if err != nil {
		sm.RUnlock()
		return errors.New("Failed Marshalling Search")
	}
	sm.RUnlock()
	err = db.SaveSearch(s.ID, s.Started, s.Private, bytes)
	if err != nil {
		return errors.New("Failed Saving Search to DB")
	}

	// Delete from Memory once saved in DB
	delete(sm.List, searchID)

	// Metrics
	metrics.SearchCount.Inc()
	metrics.SearchDuration.Observe(time.Since(start).Seconds())

	runtime.GC()

	return nil
}

// Request ...
type Request struct {
	Input   string
	Repo    string
	Private bool
	Time    time.Time
	Opts    Options
}

// Search struct auto-generated into search.pb.go

type status int

const (
	queued status = iota
	started
	completed
)
