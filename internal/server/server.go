package server

import (
	"log"
	"sync"

	"github.com/go-chi/chi"
	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/search"
	"github.com/wpdirectory/wpdir/internal/store"
)

// Server holds all the data the App needs
type Server struct {
	Store           store.DataStore
	Logger          *log.Logger
	Config          *config.Config
	Router          *chi.Mux
	PluginSearchers map[string]*search.Searcher
	ThemeSearchers  map[string]*search.Searcher
	Searches        map[string][]*index.SearchResponse
	lock            sync.RWMutex
}

// New returns a pointer to the main server struct
func New() *Server {

	return &Server{
		PluginSearchers: make(map[string]*search.Searcher),
		ThemeSearchers:  make(map[string]*search.Searcher),
		Searches:        make(map[string][]*index.SearchResponse),
	}

}

// Setup ...
func (s *Server) Setup() {

	// TODO: Pass shutdown channel for graceful shutdowns.

	// Start background Sync
	go s.startSync()

	// Start HTTP Server
	s.startHTTP()

}

// Close signals the server to gracefully shutdown.
func (s *Server) Close() {
	// Signal server to ignore new requests and finish existing.
	s.Logger.Println("Closing Server...")
}

// Shutdown will release resources and stop the server.
func (s *Server) Shutdown() {
	//s.DB.Close()
	s.Logger.Println("Server Shutdown")
}
