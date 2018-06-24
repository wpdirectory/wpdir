package server

import (
	"net/http"
)

func (s *Server) static() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Vary", "Accept-Encoding")

		w.Write([]byte(`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,shrink-to-fit=no"><meta name="theme-color" content="#000000"><link rel="manifest" href="/static/assets/manifest.json"><link rel="shortcut icon" href="/static/assets/favicon.ico"><title>WordPress Directory Searcher - WPdirectory</title><link href="/static/assets/css/main.7717625e.css" rel="stylesheet"></head><body><noscript>You need to enable JavaScript to run this app.</noscript><div id="root"></div><script type="text/javascript" src="/static/assets/js/main.216285b4.js"></script></body></html>`))
	}
}

func (s *Server) notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Vary", "Accept-Encoding")
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,shrink-to-fit=no"><meta name="theme-color" content="#000000"><link rel="manifest" href="/static/assets/manifest.json"><link rel="shortcut icon" href="/static/assets/favicon.ico"><title>WordPress Directory Searcher - WPdirectory</title><link href="/static/assets/css/main.7717625e.css" rel="stylesheet"></head><body><noscript>You need to enable JavaScript to run this app.</noscript><div id="root"></div><script type="text/javascript" src="/static/assets/js/main.216285b4.js"></script></body></html>`))
	}
}
