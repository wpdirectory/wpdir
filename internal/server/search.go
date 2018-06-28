package server

import (
	"encoding/json"
	"errors"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/theme"
	"github.com/wpdirectory/wpdir/internal/ulid"
)

type SearchManager struct {
	Queue chan string
	List  map[string]*Search
	sync.RWMutex
}

// JsonSearch is a type alias to avoid recursive json.Marshal bug
type JsonSearch Search

type Search struct {
	ID        string    `json:"id"`
	Input     string    `json:"input"`
	Repo      string    `json:"repo"`
	Matches   *Matches  `json:"matches"`
	Started   time.Time `json:"started"`
	Completed time.Time `json:"completed,omitempty"`
	Progress  int       `json:"progress"`
	Total     int       `json:"total"`
	Status    status    `json:"status"`
	Opts      Options   `json:"options"`
	Summary   *Summary  `json:"summary,omitempty"`
	Private   bool      `json:"-"`
	sync.RWMutex
}

type status int

const (
	queued status = iota
	started
	completed
)

type Matches struct {
	List  map[string][]*Match `json:"list"`
	Total int                 `json:"total,omitempty"`
	sync.RWMutex
}

type Match struct {
	Slug     string `json:"slug"`
	File     string `json:"file"`
	LineNum  int    `json:"line_num"`
	LineText string `json:"line_text"`
	//Before   []string `json:"before,omitempty"`
	//After    []string `json:"after,omitempty"`
}

type Options struct {
	CaseSensitive  bool `json:"case_sensitive"`
	LinesOfContext int  `json:"lines_context"`
}

type Summary struct {
	List  []*Item `json:"list"`
	Total int     `json:"total"`
	sync.RWMutex
}

type Item struct {
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	Version        string `json:"version"`
	Homepage       string `json:"homepage"`
	ActiveInstalls int    `json:"installs"`
	Matches        int    `json:"matches"`
}

type SearchRequest struct {
	Input   string
	Repo    string
	Private bool
	Time    time.Time
	Opts    Options
}

// Get ...
func (sm *SearchManager) Get(ID string) *Search {
	sm.Lock()
	defer sm.Unlock()
	s := sm.List[ID]
	return s
}

// Set ...
func (sm *SearchManager) Set(s *Search) {
	sm.RLock()
	defer sm.RUnlock()
	_, ok := sm.List[s.ID]
	if !ok {
		sm.List[s.ID] = s
	}
}

// Exists ...
func (sm *SearchManager) Exists(ID string) bool {
	sm.Lock()
	defer sm.Unlock()
	_, ok := sm.List[ID]
	return ok
}

// Load ...
func (sm *SearchManager) Load() int {
	i := 0
	list, err := db.GetAllFromBucket("searches")
	if err != nil {
		return i
	}

	for ID, bytes := range list {
		var s Search
		err := json.Unmarshal(bytes, &s)
		if err != nil {
			log.Printf("Failed loading search: %s %s\n", ID, err)
			db.DeleteFromBucket(ID, "searches")
		}
		sm.Set(&s)
		i++
	}

	// TODO: Order the searches before loading
	// perhaps use a temporary list to sort then
	// add to the SearchManager

	return i
}

// Empty ...
func (sm *SearchManager) Empty() int {
	i := 0
	list, err := db.GetAllFromBucket("searches")
	if err != nil {
		return i
	}

	for ID := range list {
		db.DeleteFromBucket(ID, "searches")
		i++
	}

	return i
}

// NewSearch ...
func (sm *SearchManager) NewSearch(sr SearchRequest) string {
	sm.RLock()
	defer sm.RUnlock()

	ID := ulid.New()
	sm.List[ID] = &Search{
		ID:      ID,
		Input:   sr.Input,
		Repo:    sr.Repo,
		Private: sr.Private,
		Opts:    sr.Opts,
		Matches: &Matches{
			List:  make(map[string][]*Match),
			Total: 0,
		},
		Status: queued,
	}

	sm.Queue <- ID

	return ID
}

// SearchWorker ...
func (s *Server) SearchWorker() {

	for {
		searchID := <-s.Searches.Queue
		err := s.processSearch(searchID)
		if err != nil {
			log.Printf("Searched failed: %s\n", err)
		}
	}

}

// processSearch ...
func (s *Server) processSearch(ID string) error {
	s.Searches.Lock()
	srch := s.Searches.List[ID]
	s.Searches.Unlock()

	var totalMatches int
	srch.Started = time.Now()
	srch.Status = started

	sum := &Summary{
		List:  []*Item{},
		Total: 0,
	}

	opts := &index.SearchOptions{
		Offset:         0,
		Limit:          0,
		LinesOfContext: 2,
		IgnoreCase:     false,
	}

	limit := runtime.NumCPU() - 2
	limiter := make(chan struct{}, limit)

	switch srch.Repo {
	case "plugins":
		pr := s.Plugins.(*repo.PluginRepo)
		srch.Total = pr.Len()

		for _, p := range pr.List {
			if !p.HasIndex() || p.Status != 0 {
				continue
			}
			limiter <- struct{}{}

			go func(p *plugin.Plugin) {
				resp, err := p.Searcher.Search(srch.Input, p.Slug, opts)
				if err != nil {
					<-limiter
					return
				}
				if len(resp.Matches) == 0 {
					<-limiter
					return
				}
				item := &Item{
					Slug: p.Slug,
				}
				srch.Matches.RLock()
				for _, result := range resp.Matches {
					for _, match := range result.Matches {
						totalMatches++
						sum.Total++
						item.Matches++
						m := &Match{
							Slug:     p.Slug,
							File:     result.Filename,
							LineNum:  match.LineNumber,
							LineText: match.Line,
						}
						// Does this need Locks?
						srch.Matches.Total++
						srch.Matches.List[p.Slug] = append(srch.Matches.List[p.Slug], m)
					}
				}
				srch.Matches.RUnlock()

				sum.RLock()
				sum.List = append(sum.List, item)
				sum.RUnlock()
				<-limiter
			}(p)
			srch.RLock()
			srch.Progress++

			sum.Lock()
			srch.Summary = sum
			sum.Unlock()

			srch.RUnlock()
		}

		break

	case "themes":
		tr := s.Themes.(*repo.ThemeRepo)
		srch.Total = tr.Len()

		for _, t := range tr.List {
			srch.Progress++
			if !t.HasIndex() || t.Status != 0 {
				continue
			}
			limiter <- struct{}{}

			go func(t *theme.Theme) {
				resp, err := t.Searcher.Search(srch.Input, t.Slug, opts)
				if err != nil {
					<-limiter
					return
				}
				if len(resp.Matches) == 0 {
					<-limiter
					return
				}

				item := &Item{
					Slug: t.Slug,
				}
				for _, result := range resp.Matches {
					for _, match := range result.Matches {
						totalMatches++
						sum.Total++
						item.Matches++
						m := &Match{
							Slug:     t.Slug,
							File:     result.Filename,
							LineNum:  match.LineNumber,
							LineText: match.Line,
						}
						srch.Matches.RLock()
						srch.Matches.Total++
						srch.Matches.List[t.Slug] = append(srch.Matches.List[t.Slug], m)
						srch.Matches.RUnlock()
					}
				}
				sum.RLock()
				sum.List = append(sum.List, item)
				sum.RUnlock()
				<-limiter
			}(t)
			srch.Lock()
			srch.Summary = sum
			srch.Unlock()
		}

		break

	default:
		return errors.New("Not a valid respository name")
	}

	srch.Completed = time.Now()
	srch.Status = completed

	bytes, err := json.Marshal(srch)
	if err != nil {
		return err
	}

	err = db.PutToBucket(srch.ID, bytes, "searches")
	if err != nil {
		log.Printf("Could not save search to DB: %s\n", err)
	}

	return nil
}

// MarshalJSON handles locking Search during json.Marshal
func (srch *Search) MarshalJSON() ([]byte, error) {
	srch.RLock()
	defer srch.RUnlock()

	return json.Marshal(JsonSearch(*srch))
}
