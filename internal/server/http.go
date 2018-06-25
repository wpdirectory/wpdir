package server

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// startHTTP starts the HTTP server.
func (s *Server) startHTTP() {

	s.Router = chi.NewRouter()

	// Middleware Stack
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.DefaultCompress)
	s.Router.Use(middleware.RedirectSlashes)

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

	filesDir := filepath.Join("D:/projects/go/src/github.com/wpdirectory/wpdir", "web", "build")
	FileServer(s.Router, "/static", http.Dir(filesDir))

	s.routes()

	srvHTTP := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + req.Host + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
		Addr: ":" + s.Config.Ports.HTTP,
	}
	go func() { log.Fatal(srvHTTP.ListenAndServe()) }()

	srvHTTPS := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig(),
		Handler:      s.Router,
		Addr:         ":" + s.Config.Ports.HTTPS,
	}

	cert := filepath.Join(s.Config.WD, "certs", "wpdirectory.net.crt")
	key := filepath.Join(s.Config.WD, "certs", "wpdirectory.net.key")

	log.Fatal(srvHTTPS.ListenAndServeTLS(cert, key))

}

func (s *Server) routes() {

	// Add Routes
	s.Router.Get("/", s.static())
	s.Router.Get("/search/{id}", s.static())
	s.Router.Get("/searches", s.static())
	s.Router.Get("/repos", s.static())
	s.Router.Get("/about", s.static())

	// Need to disable RedirectSlashes middleware to enable this
	// redirects to /debug/prof/ which causes redirect loop
	//s.Router.Mount("/debug", middleware.Profiler())

	// Add API v1 routes
	s.Router.Mount("/api/v1", s.apiRoutes())

	// Handle NotFound
	s.Router.NotFound(s.notFound())

}

func (s *Server) apiRoutes() chi.Router {

	r := chi.NewRouter()

	r.Get("/search/{id}", s.getSearch())
	r.Post("/search/new", s.createSearch())
	r.Get("/searches/{limit}", s.getSearches())
	r.Get("/search/matches/{id}/{slug}", s.getSearchMatches())

	r.Post("/file", s.getMatchFile())

	r.Get("/repo/{name}", s.getRepo())
	r.Get("/repos/overview", s.getRepoOverview())

	r.Get("/plugin/{slug}", s.getPlugin())

	r.Get("/theme/{slug}", s.getTheme())

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
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
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

// Setup the TLS/SSL configuration.
func tlsConfig() *tls.Config {

	return &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

}
