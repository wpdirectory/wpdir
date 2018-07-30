package repo

import (
	"sync"

	"github.com/wpdirectory/wpdir/internal/filestats"
	"github.com/wpdirectory/wpdir/internal/index"
)

// Extension holds data about a Plugin or Theme
type Extension struct {
	Slug                   string     `json:"slug"`
	Name                   string     `json:"name,omitempty"`
	Version                string     `json:"version,omitempty"`
	Author                 string     `json:"author,omitempty"`
	AuthorProfile          string     `json:"author_profile,omitempty"`
	Contributors           [][]string `json:"contributors,omitempty"`
	Requires               string     `json:"requires,omitempty"`
	Tested                 string     `json:"tested,omitempty"`
	RequiresPHP            string     `json:"requires_php,omitempty"`
	Rating                 int        `json:"rating,omitempty"`
	Ratings                []Rating   `json:"ratings,omitempty"`
	NumRatings             int        `json:"num_ratings,omitempty"`
	SupportThreads         int        `json:"support_threads,omitempty"`
	SupportThreadsResolved int        `json:"support_threads_resolved,omitempty"`
	ActiveInstalls         int        `json:"active_installs,omitempty"`
	Downloaded             int        `json:"downloaded,omitempty"`
	LastUpdated            string     `json:"last_updated,omitempty"`
	Added                  string     `json:"added,omitempty"`
	Homepage               string     `json:"homepage,omitempty"`
	Sections               struct {
		Description string `json:"description,omitempty"`
		Faq         string `json:"faq,omitempty"`
		Changelog   string `json:"changelog,omitempty"`
		Screenshots string `json:"screenshots,omitempty"`
	} `json:"sections,omitempty"`
	ShortDescription string       `json:"short_description,omitempty"`
	DownloadLink     string       `json:"download_link,omitempty"`
	Screenshots      []Screenshot `json:"screenshots,omitempty"`
	Tags             [][]string   `json:"tags,omitempty"`
	StableTag        string       `json:"stable_tag,omitempty"`
	Versions         [][]string   `json:"versions,omitempty"`
	DonateLink       string       `json:"donate_link,omitempty"`
	Status           status       `json:"status,omitempty"`
	index            *index.Index
	Stats            *filestats.Stats `json:"stats,omitempty"`
	sync.RWMutex
}

// Rating contains information about ratings of a specific star level (0-5)
type Rating struct {
	Stars  string `json:"stars"`
	Number int    `json:"number"`
}

// Screenshot contains the source and caption of a screenshot
type Screenshot struct {
	Src     string `json:"src"`
	Caption string `json:"caption"`
}

type status int

const (
	// Open shows we have files and API info stored
	Open status = iota
	// Closed shows we cannot get data
	Closed
)

// NewExt returns a new Extension.
func NewExt(slug string) *Extension {
	return &Extension{
		Slug:   slug,
		Status: Closed,
	}
}

// GetStatus returns the Status as a string
func (e *Extension) GetStatus() string {
	e.RLock()
	defer e.RUnlock()

	switch e.Status {
	case Open:
		return "Open"
	case Closed:
		return "Closed"
	default:
		return "Invalid Status"
	}
}

// SetStatus sets the Extension Status
func (e *Extension) SetStatus(s status) {
	e.Lock()
	defer e.Unlock()
	e.Status = s
}

// SwapIndexes ...
func (e *Extension) SwapIndexes(idx *index.Index) error {
	e.Lock()
	defer e.Unlock()

	oldIdx := e.index
	e.index = idx

	if oldIdx != nil {
		return oldIdx.Destroy()
	}

	return nil
}

// Dir returns the index dir
func (e *Extension) Dir() string {
	e.index.RLock()
	defer e.index.RUnlock()

	return e.index.Ref.Dir()
}

// Search performs a basic search on the current index using the supplied pattern
// and the options.
func (e *Extension) Search(pat, slug string, opt *index.SearchOptions) (*index.SearchResponse, error) {
	e.RLock()
	defer e.RUnlock()
	return e.index.Search(pat, slug, opt)
}
