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
	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/search/queue"
	"github.com/wpdirectory/wpdir/internal/theme"
	"github.com/wpdirectory/wpdir/internal/ulid"
)

// Manager controls the processing and storage of searches
type Manager struct {
	Queue   *queue.Queue
	List    map[string]*Search
	Plugins *repo.PluginRepo
	Themes  *repo.ThemeRepo
	sync.RWMutex
}

// NewManager returns a new SearchManager struct
func NewManager() *Manager {
	return &Manager{
		Queue: queue.New(100),
		List:  make(map[string]*Search),
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

type SummaryList struct {
	List  map[string]*Result
	Total uint64
	sync.RWMutex
}

type MatchList struct {
	List map[string]*Matches
	sync.RWMutex
}

// processSearch ...
func (sm *Manager) processSearch(ID string) error {
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

	limit := runtime.NumCPU() - 2
	limiter := make(chan struct{}, limit)

	var wg sync.WaitGroup

	switch srch.Repo {
	case "plugins":
		pr := sm.Plugins
		total = pr.Len()
		for _, p := range pr.List {
			current++
			sm.Lock()
			srch.Progress = uint32(math.Round((float64(current) / float64(total)) * 100.00))
			srch.Matches = uint32(totalMatches)
			sm.Unlock()
			if !p.HasIndex() || p.Status != 0 {
				continue
			}
			wg.Add(1)
			limiter <- struct{}{}

			go func(p *plugin.Plugin, input string, sum *SummaryList, matchlist *MatchList, totalMatches *uint64, wg *sync.WaitGroup) {
				p.RLock()
				defer p.RUnlock()
				resp, err := p.Searcher.Search(input, p.Slug, opts)
				if err != nil || len(resp.Matches) == 0 {
					wg.Done()
					<-limiter
					return
				}
				var pMatches uint64
				for i := 0; i < len(resp.Matches); i++ {
					if resp.Matches[i].Matches == nil {
						continue
					}
					pMatches = uint64(len(resp.Matches[i].Matches))
					atomic.AddUint64(totalMatches, pMatches)
					ms := &Matches{}
					for j := 0; j < len(resp.Matches[i].Matches); j++ {
						m := &Match{
							Slug:     p.Slug,
							File:     resp.Matches[i].Filename,
							LineNum:  uint32(resp.Matches[i].Matches[j].LineNumber),
							LineText: resp.Matches[i].Matches[j].Line,
						}
						ms.List = append(ms.List, m)
					}
					matchList.Lock()
					matchList.List[p.Slug] = ms
					matchList.Unlock()
				}
				r := &Result{
					Slug:           p.Slug,
					Name:           p.Name,
					Version:        p.Version,
					Homepage:       p.Homepage,
					ActiveInstalls: uint32(p.ActiveInstalls),
					Matches:        uint32(pMatches),
				}
				sum.Lock()
				sum.List[p.Slug] = r
				sum.Unlock()
				wg.Done()
				<-limiter
			}(p, input, sum, matchList, &totalMatches, &wg)
		}

		break
	case "themes":
		tr := sm.Themes
		total = tr.Len()
		for _, t := range tr.List {
			current++
			sm.Lock()
			srch.Progress = uint32(math.Round((float64(current) / float64(total)) * 100.00))
			srch.Matches = uint32(totalMatches)
			sm.Unlock()
			if !t.HasIndex() || t.Status != 0 {
				continue
			}
			wg.Add(1)
			limiter <- struct{}{}

			go func(t *theme.Theme, input string, sum *SummaryList, matchlist *MatchList, totalMatches *uint64, wg *sync.WaitGroup) {
				t.RLock()
				defer t.RUnlock()
				resp, err := t.Searcher.Search(input, t.Slug, opts)
				if err != nil || len(resp.Matches) == 0 {
					wg.Done()
					<-limiter
					return
				}
				var tMatches uint64
				for i := 0; i < len(resp.Matches); i++ {
					if resp.Matches[i].Matches == nil {
						continue
					}
					tMatches = uint64(len(resp.Matches[i].Matches))
					atomic.AddUint64(totalMatches, tMatches)
					ms := &Matches{}
					for j := 0; j < len(resp.Matches[i].Matches); j++ {
						m := &Match{
							Slug:     t.Slug,
							File:     resp.Matches[i].Filename,
							LineNum:  uint32(resp.Matches[i].Matches[j].LineNumber),
							LineText: resp.Matches[i].Matches[j].Line,
						}
						ms.List = append(ms.List, m)
					}
					matchList.Lock()
					matchList.List[t.Slug] = ms
					matchList.Unlock()
				}
				r := &Result{
					Slug:           t.Slug,
					Name:           t.Name,
					Version:        t.Version,
					Homepage:       t.Homepage,
					ActiveInstalls: uint32(t.ActiveInstalls),
					Matches:        uint32(tMatches),
				}
				sum.Lock()
				sum.List[t.Slug] = r
				sum.Unlock()
				wg.Done()
				<-limiter
			}(t, input, sum, matchList, &totalMatches, &wg)
		}

		break
	default:
		return errors.New("Not a valid respository name")
	}

	wg.Wait()

	summary := &Summary{
		List:  make(map[string]*Result),
		Total: sum.Total,
	}
	for key, result := range sum.List {
		summary.List[key] = result
	}

	// TODO: Store Search, Summary and MatchList in DB.
	sm.RLock()
	// Copy Search so we can close the lock earlier
	s := *srch
	bytes, err := srch.Marshal()
	if err != nil {
		return errors.New("Failed Marshalling Search")
	}
	sm.RUnlock()
	err = db.SaveSearch(s.ID, s.Started, s.Private, bytes)
	if err != nil {
		return errors.New("Failed Saving Search to DB")
	}

	bytes, err = summary.Marshal()
	if err != nil {
		return errors.New("Failed Marshalling Summary")
	}
	err = db.SaveSummary(s.ID, bytes)
	if err != nil {
		return errors.New("Failed Saving Summary to DB")
	}

	for slug, matches := range matchList.List {
		bytes, err = matches.Marshal()
		if err != nil {
			return errors.New("Failed Marshalling Matches")
		}
		err = db.SaveMatches(s.ID, slug, bytes)
		if err != nil {
			return errors.New("Failed Saving Summary to DB")
		}
	}

	sm.Lock()
	srch.Completed = time.Now().Format(time.RFC3339)
	srch.Status = Completed
	srch.Matches = uint32(totalMatches)
	sm.Unlock()

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
