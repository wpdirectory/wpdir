package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/search"
)

// Server holds all the data the App needs
type Server struct {
	Logger  *log.Logger
	Config  *config.Config
	Router  *chi.Mux
	Manager *search.Manager
	http    *http.Server
	https   *http.Server
}

// New returns a pointer to the main server struct
func New(log *log.Logger, config *config.Config, fresh *bool) *Server {

	// Init Repos
	pr := repo.New(config, log, "plugins", 0)
	tr := repo.New(config, log, "themes", 0)

	// TODO: Auto generate in background
	//pr.GenerateInstallsChart()
	//tr.GenerateInstallsChart()
	//pr.GenerateSizeChart()
	//tr.GenerateSizeChart()

	sm := search.NewManager(config.SearchWorkers)
	sm.Plugins = pr
	sm.Themes = tr

	// Debug Delete Searches
	// Need to reset after break code changes
	//sm.Empty()

	s := &Server{
		Config:  config,
		Logger:  log,
		Manager: sm,
	}

	// Load Existing Data
	go s.LoadData(fresh)

	return s
}

// LoadData loads all existing DB and Index data
func (s *Server) LoadData(fresh *bool) {
	// Load Existing from DB
	s.Manager.Plugins.LoadExisting()
	s.Manager.Themes.LoadExisting()

	// Initial List
	err := s.Manager.Plugins.UpdateList(fresh)
	if err != nil {
		s.Logger.Fatalf("Could not get initial plugin list")
	}
	err = s.Manager.Themes.UpdateList(fresh)
	if err != nil {
		s.Logger.Fatalf("Could not get initial theme list")
	}

	// Start Update Workers
	// These process updates from the queue
	repo.StartUpdateWorkers(s.Config.UpdateWorkers, s.Manager.Plugins, s.Manager.Themes)

	// Start Worker to Process Searches
	go s.Manager.Worker()

	s.Manager.Lock()
	s.Manager.Loaded = true
	s.Manager.Unlock()
}

// Setup starts the HTTP Server
func (s *Server) Setup() {
	s.startUp()
}

// Shutdown will release resources and stop the server.
func (s *Server) Shutdown(ctx context.Context) {
	s.http.Shutdown(ctx)
	s.https.Shutdown(ctx)
}
