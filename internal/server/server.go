package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/store"
)

var wd string

// Server holds all the data the App needs
type Server struct {
	Store  store.DataStore
	Logger *log.Logger
	Config *config.Config
	Email  string
	Router *chi.Mux
}

// New returns a pointer to the main server struct
func New() *Server {

	return &Server{}

}

// Setup ...
func (s *Server) Setup() {

	wd = s.Config.WD

	s.Router = chi.NewRouter()

	// Middleware Stack
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.DefaultCompress)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	s.Router.Use(middleware.Timeout(60 * time.Second))

	// TODO: Remove this for prod?
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	s.Router.Use(cors.Handler)

	wd, _ := os.Getwd()
	filesDir := filepath.Join(wd, "assets")
	FileServer(s.Router, "/static", http.Dir(filesDir))

	s.routes()

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      s.Router,
		Addr:         ":" + s.Config.HTTP.Port,
	}

	log.Fatal(srv.ListenAndServe())

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

func (s *Server) routes() {

	// Add Routes
	s.Router.Get("/", s.static())
	s.Router.Get("/search/", s.static())
	s.Router.Get("/search/{id}/", s.static())
	s.Router.Get("/searches/", s.static())
	s.Router.Get("/status/", s.static())
	s.Router.Get("/about/", s.static())

	// Add API v1 routes
	s.Router.Mount("/api/v1", s.apiRoutes())

	// Handle NotFound
	s.Router.NotFound(s.notFound())

}

func (s *Server) apiRoutes() chi.Router {

	r := chi.NewRouter()

	r.Get("/search/{id}/", s.getSearch())
	r.Post("/search/new/", s.createSearch())
	r.Get("/searches/", s.getSearchList())

	return r

}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		fs.ServeHTTP(w, r)
	}))
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Panicf("Failed to encode JSON: %v\n", err)
	}
}

func writeResp(w http.ResponseWriter, data interface{}) {
	writeJSON(w, data, http.StatusOK)
}

func writeError(w http.ResponseWriter, err error, status int) {
	writeJSON(w, map[string]string{
		"Error": err.Error(),
	}, status)
}
