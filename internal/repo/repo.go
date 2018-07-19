package repo

import (
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/theme"
	"github.com/wpdirectory/wporg"
)

var (
	httpClient *http.Client
)

// Repo ...
type Repo interface {
	Len() uint64
	Rev() int

	Exists(slug string) bool
	Get(slug string) Extension
	Add(slug string)
	Remove(slug string)
	UpdateIndex(idx *index.Index) error
	UpdateList() error

	load() error
	save() error
	LoadExisting()

	QueueUpdate(slug string)
	UpdateWorker()
	StartWorkers()
	ProcessUpdate(slug string) error

	Summary() *Summary
}

// New returns a new Repo
func New(t string, c *config.Config, l *log.Logger) Repo {
	// Setup HTTP Client
	opt := func(c *wporg.Client) {
		c.HTTPClient = httpClient
	}
	api := wporg.NewClient(opt)
	var repo Repo
	switch t {
	case "plugins":
		repo = &PluginRepo{
			cfg:         c,
			log:         l,
			api:         api,
			List:        make(map[string]*plugin.Plugin),
			Revision:    1904883,
			UpdateQueue: make(chan string, 100000),
		}
	case "themes":
		repo = &ThemeRepo{
			cfg:         c,
			log:         l,
			api:         api,
			List:        make(map[string]*theme.Theme),
			Revision:    96064,
			UpdateQueue: make(chan string, 75000),
		}
	}
	// Load Existing Data
	err := repo.load()
	if err != nil {
		l.Printf("Repo (%s) could not load data: %s\n", t, err)
	}

	return repo
}

// Summary ...
type Summary struct {
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

func init() {
	var netTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}

	httpClient = &http.Client{
		Timeout:   time.Second * time.Duration(120),
		Transport: netTransport,
	}
}
