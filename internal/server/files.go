package server

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/wpdirectory/wpdir/internal/repo"
)

func (s *Server) getFilePath(repository, slug, file string) (string, error) {

	// Protect against directory traversal attacks
	if containsDotDot(repository) || containsDotDot(slug) || containsDotDot(file) {
		return "", errors.New("Paths must not include '..'")
	}

	switch repository {
	case "plugins":
		if !s.Manager.Plugins.Exists(slug) {
			return "", errors.New("No matching plugin")
		}
		p := s.Manager.Plugins.Get(slug)
		if p.Status != repo.Open {
			return "", errors.New("Plugin is Closed")
		}

		p.RLock()
		dir := p.Dir()
		p.RUnlock()

		path := filepath.Join(dir, "raw", file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return "", errors.New("File not found")
		}

		return path, nil

	case "themes":
		if !s.Manager.Themes.Exists(slug) {
			return "", errors.New("No matching theme")
		}

		t := s.Manager.Themes.Get(slug)
		if t.Status != repo.Open {
			return "", errors.New("Theme has no indexed files")
		}

		t.RLock()
		dir := t.Dir()
		t.RUnlock()

		path := filepath.Join(dir, "raw", file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return "", errors.New("File not found")
		}

		return path, nil

	default:
		return "", errors.New("No matching repository")
	}
}

// Coped from Go's net/http package
func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }
