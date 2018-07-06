package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/repo"
)

// Server holds all the data the App needs
type Server struct {
	Logger   *log.Logger
	Config   *config.Config
	Router   *chi.Mux
	Plugins  repo.Repo
	Themes   repo.Repo
	Searches *SearchManager
	http     *http.Server
	https    *http.Server
}

// New returns a pointer to the main server struct
func New(log *log.Logger, config *config.Config) *Server {

	// Init Repos
	pr := repo.New("plugins", config, log)
	tr := repo.New("themes", config, log)

	// Load Existing from DB
	pr.LoadExisting()
	tr.LoadExisting()

	// Initial List
	err := pr.UpdateList()
	if err != nil {
		log.Fatalf("Could not get initial plugin list")
	}
	err = tr.UpdateList()
	if err != nil {
		log.Fatalf("Could not get initial theme list")
	}

	sm := &SearchManager{
		Queue: make(chan string, 200),
		List:  make(map[string]*Search),
	}

	// Debug Delete Searches
	// Need to reset after break code changes
	//del := sm.Empty()
	//log.Printf("Deleted %d searches.", del)

	// Load Existing Searches
	count := sm.Load()
	log.Printf("Loaded %d searches.", count)

	s := &Server{
		Config:   config,
		Logger:   log,
		Plugins:  pr.(*repo.PluginRepo),
		Themes:   tr.(*repo.ThemeRepo),
		Searches: sm,
	}

	// Start Workers
	go pr.UpdateWorker()
	go tr.UpdateWorker()

	// Start Worker to Process Searches
	go s.SearchWorker()

	return s
}

// Setup ...
func (s *Server) Setup() {

	// TODO: Pass shutdown channel for graceful shutdowns.

	// Start HTTP Server
	s.startHTTP()

}

// Close signals the server to gracefully shutdown.
func (s *Server) Close() {
	// Signal server to ignore new requests and finish existing.
	s.Logger.Println("Closing Server...")
}

// Shutdown will release resources and stop the server.
func (s *Server) Shutdown(ctx context.Context) {
	//s.DB.Close()
	s.Logger.Println("Server Shutdown")
}
