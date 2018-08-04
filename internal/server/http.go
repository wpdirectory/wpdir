package server

import (
	"crypto/tls"
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
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"github.com/wpdirectory/wpdir/internal/data"
	"github.com/wpdirectory/wpdir/internal/limit"
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
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	s.Router.Use(cors.Handler)

	FileServer(s.Router, "/assets")

	s.routes()

	cert := filepath.Join(s.Config.WD, "certs", "wpdirectory.net.crt")
	key := filepath.Join(s.Config.WD, "certs", "wpdirectory.net.key")

	// Serve HTTP if no cert found, HTTPS otherwise
	if _, err := os.Stat(cert); os.IsNotExist(err) {
		s.Logger.Println("No certs found, serving on HTTP Port")
		s.serveHTTP()
	} else {
		s.Logger.Println("Found certs, serving on HTTP and HTTPS Ports")
		s.serveHTTPS(cert, key)
	}
}

func (s *Server) serveHTTP() {
	s.http = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      s.Router,
		Addr:         ":" + s.Config.Ports.HTTP,
	}
	go func() { log.Fatal(s.http.ListenAndServe()) }()
}

func (s *Server) serveHTTPS(cert, key string) {
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	s.http = &http.Server{
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
	go func() { log.Fatal(s.http.ListenAndServe()) }()

	s.https = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      s.Router,
		Addr:         ":" + s.Config.Ports.HTTPS,
	}
	go func() { log.Fatal(s.https.ListenAndServeTLS(cert, key)) }()
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

	middleware := stdlib.NewMiddleware(limit.New(), stdlib.WithForwardHeader(true))

	r.Get("/search/{id}", s.getSearch())
	r.Post("/search/new", middleware.Handler(s.createSearch()).(http.HandlerFunc))
	r.Get("/searches/{limit}", s.getSearches())
	r.Get("/search/matches/{id}/{slug}", s.getSearchMatches())

	r.Get("/search/summary/{id}", s.getSearchSummary())

	r.Post("/file", s.getMatchFile())

	r.Get("/repo/{name}", s.getRepo())
	r.Get("/repos/overview", s.getRepoOverview())

	r.Get("/plugin/{slug}", s.getPlugin())

	r.Get("/theme/{slug}", s.getTheme())

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.FileServer(data.Assets)

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
