package searcher

import (
	"sync"

	"github.com/wpdirectory/wpdir/internal/index"
)

type Searcher struct {
	idx *index.Index
	sync.RWMutex
}

// New ...
func New(ref *index.IndexRef) (*Searcher, error) {

	idx, err := ref.Open()
	if err != nil {
		return &Searcher{}, err
	}

	s := &Searcher{
		idx: idx,
	}

	return s, nil
}

// SwapIndexes performs atomic swap of index in the searcher so that the new
// index is made "live".
func (s *Searcher) SwapIndexes(idx *index.Index) error {
	s.Lock()
	defer s.Unlock()

	oldIdx := s.idx
	s.idx = idx

	if oldIdx != nil {
		return oldIdx.Destroy()
	}

	return nil
}

// Dir returns the index dir
func (s *Searcher) Dir() string {
	s.Lock()
	defer s.Unlock()

	return s.idx.Ref.Dir()
}

// Perform a basic search on the current index using the supplied pattern
// and the options.
//
// TODO(knorton): pat should really just be a part of SearchOptions
func (s *Searcher) Search(pat, slug string, opt *index.SearchOptions) (*index.SearchResponse, error) {
	s.RLock()
	defer s.RUnlock()
	return s.idx.Search(pat, slug, opt)
}
