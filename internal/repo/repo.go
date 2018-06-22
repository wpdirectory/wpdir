package repo

import (
	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/theme"
)

// Repo ...
type Repo interface {
	Len() int
	Rev() int

	Exists(slug string) bool
	Get(slug string) Extension
	Add(slug string)
	Remove(slug string)
	UpdateIndex(idx *index.Index) error
	UpdateList() error

	LoadExisting()

	QueueUpdate(slug string)
	UpdateWorker()
	ProcessUpdate(slug string) error

	Summary() *RepoSummary
}

// New returns a new Repo
func New(t string, c *config.Config) Repo {
	var repo Repo
	switch t {
	case "plugins":
		repo = &PluginRepo{
			Config:      c,
			List:        make(map[string]*plugin.Plugin),
			Revision:    0,
			UpdateQueue: make(chan string, 100000),
		}
	case "themes":
		repo = &ThemeRepo{
			Config:      c,
			List:        make(map[string]*theme.Theme),
			Revision:    0,
			UpdateQueue: make(chan string, 100000),
		}
	}
	return repo
}

// RepoSummary ...
type RepoSummary struct {
	Revision int `json:"revision"`
	Total    int `json:"total"`
	Closed   int `json:"closed"`
	Queue    int `json:"queue"`
}

// Extension ...
type Extension interface {
	GetStatus() string
	HasIndex() bool
	SetIndexed(idx bool)
	LoadAPIData() error
	Update() error
	Save() error
}
