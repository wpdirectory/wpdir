package server

import (
	"net/http"
)

var (
	html = []byte(`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,shrink-to-fit=no"><meta name="theme-color" content="#000000"><meta property="og:locale" content="en_US"/><meta property="og:type" content="website"/><meta property="og:title" content="WordPress Directory Searcher - WPDirectory"/><meta property="og:url" content="https://wpdirectory.net/"/><meta property="og:site_name" content="WPDirectory"/><meta name="description" content="Lightning fast regex searching of code in the WordPress Plugin and Theme Directories. Start searching now!"/><meta name="keywords" content="search, regex, wordpress, plugins, themes"/><meta name="author" content="Peter Booker"/><meta name="contact" content="mail@peterbooker.com"/><link rel="manifest" href="/static/assets/manifest.json"><link rel="shortcut icon" href="/static/assets/favicon.ico"><title>WordPress Directory Searcher - WPDirectory</title><link href="/static/assets/css/main.40ec9d91.css" rel="stylesheet"></head><body><noscript>You need to enable JavaScript to run this app.</noscript><div id="root"></div><script type="text/javascript" src="/static/assets/js/main.821b61d4.js"></script></body></html>`)
)

func (s *Server) static() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Vary", "Accept-Encoding")

		w.Write(html)
	}
}

func (s *Server) notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Vary", "Accept-Encoding")
		w.WriteHeader(http.StatusNotFound)

		w.Write(html)
	}
}
