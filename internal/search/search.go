package search

import (
	"sync"

	"github.com/wpdirectory/wpdir/internal/index"
)

type Searcher struct {
	Ext *Extension
	idx *index.Index
	lck sync.RWMutex
}

type Extension struct {
	Name     string
	Slug     string
	Dir      string
	Revision string
}

// New ...
func New(name string, slug string, rev string, ref *index.IndexRef) (*Searcher, error) {

	extn := &Extension{
		Name:     slug,
		Slug:     slug,
		Dir:      "plugins",
		Revision: rev,
	}

	idx, err := ref.Open()
	if err != nil {
		return &Searcher{}, err
	}

	s := &Searcher{
		idx: idx,
		Ext: extn,
	}

	return s, nil
}

// SwapIndexes performs atomic swap of index in the searcher so that the new
// index is made "live".
func (s *Searcher) SwapIndexes(idx *index.Index) error {
	s.lck.Lock()
	defer s.lck.Unlock()

	oldIdx := s.idx
	s.idx = idx

	return oldIdx.Destroy()
}

// Perform a basic search on the current index using the supplied pattern
// and the options.
//
// TODO(knorton): pat should really just be a part of SearchOptions
func (s *Searcher) Search(pat string, opt *index.SearchOptions) (*index.SearchResponse, error) {
	s.lck.RLock()
	defer s.lck.RUnlock()
	return s.idx.Search(pat, opt)
}
