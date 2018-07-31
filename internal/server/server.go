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

	// Load Existing from DB
	pr.LoadExisting()
	tr.LoadExisting()

	// Initial List
	err := pr.UpdateList(fresh)
	if err != nil {
		log.Fatalf("Could not get initial plugin list")
	}
	err = tr.UpdateList(fresh)
	if err != nil {
		log.Fatalf("Could not get initial theme list")
	}

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

	// Start Update Workers
	go repo.StartUpdateWorkers(config.UpdateWorkers, pr, tr)

	// Start Worker to Process Searches
	go sm.Worker()

	return s
}

// Setup ...
func (s *Server) Setup() {

	// TODO: Pass shutdown channel for graceful shutdowns.

	// Start HTTP Server
	s.startHTTP()

}

// Shutdown will release resources and stop the server.
func (s *Server) Shutdown(ctx context.Context) {
	s.http.Shutdown(ctx)
	s.https.Shutdown(ctx)
	s.Logger.Println("Server Shutdown")
}
