package server

import (
	"net/http"
)

func (s *Server) static() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Add("Vary", "Accept-Encoding")

		w.Write([]byte(`Index`))

		//filesDir := filepath.Join(wd, "assets")
		//http.ServeFile(w, r, path.Join(filesDir, "/index.html"))
	}
}

func (s *Server) notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Add("Vary", "Accept-Encoding")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`Not Found`))

		//filesDir := filepath.Join(wd, "assets")
		//http.ServeFile(w, r, path.Join(filesDir, "/index.html"))
	}
}
