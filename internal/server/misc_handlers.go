package server

import (
	"io"
	"net/http"
	"strings"

	"github.com/wpdirectory/wpdir/internal/data"
)

// TODO: Rewrite how the index.html file is served
// Need to embed SEO related info too

func (s *Server) static() http.HandlerFunc {
	file, err := data.Assets.Open("/index.html")
	if err != nil {
		s.Logger.Fatalf("Failed to open index.html: %s\n", err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		s.Logger.Fatalf("Failed to get file info index.html: %s\n", err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		s.Logger.Fatalf("Failed to read index.html: %s\n", err)
	}

	html := s.addConfig(buffer)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Vary", "Accept-Encoding")

		w.Write(html)
	}
}

func (s *Server) notFound() http.HandlerFunc {
	file, err := data.Assets.Open("/index.html")
	if err != nil {
		s.Logger.Fatalf("Failed to open index.html: %s\n", err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		s.Logger.Fatalf("Failed to get file info index.html: %s\n", err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		s.Logger.Fatalf("Failed to read index.html: %s\n", err)
	}

	html := s.addConfig(buffer)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Vary", "Accept-Encoding")
		w.WriteHeader(http.StatusNotFound)

		w.Write(html)
	}
}

func (s *Server) addConfig(html []byte) []byte {
	// Embed Version into HTML
	html = []byte(strings.Replace(string(html), "%VERSION%", s.Config.Version, 1))

	// Embed Commit Hash into HTML
	html = []byte(strings.Replace(string(html), "%COMMIT%", s.Config.Commit, 1))

	// Embed Date into HTML
	html = []byte(strings.Replace(string(html), "%DATE%", s.Config.Date, 1))

	// Embed Hostname into HTML, remove trailing slash
	host := strings.TrimRight(s.Config.Host, "/")
	html = []byte(strings.Replace(string(html), "%HOSTNAME%", host, 1))

	// Embed HTTP Config into HTML
	html = []byte(strings.Replace(string(html), "%TIMEOUT%", "5000", 1))

	return html
}
