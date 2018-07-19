package server

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/theme"
)

func (s *Server) getFilePath(repo, slug, file string) (string, error) {

	// Protect against directory traversal attacks
	if containsDotDot(repo) || containsDotDot(slug) || containsDotDot(file) {
		return "", errors.New("Paths must not include '..'")
	}

	switch repo {
	case "plugins":
		if !s.Manager.Plugins.Exists(slug) {
			return "", errors.New("No matching plugin")
		}
		p := s.Manager.Plugins.Get(slug).(*plugin.Plugin)
		if !p.HasIndex() {
			return "", errors.New("Plugin has no indexed files")
		}

		p.Searcher.Lock.RLock()
		dir := p.Searcher.Dir()
		p.Searcher.Lock.RUnlock()

		path := filepath.Join(dir, "raw", file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return "", errors.New("File not found")
		}

		return path, nil

	case "themes":
		if !s.Manager.Themes.Exists(slug) {
			return "", errors.New("No matching theme")
		}

		t := s.Manager.Themes.Get(slug).(*theme.Theme)
		if !t.HasIndex() {
			return "", errors.New("Theme has no indexed files")
		}

		t.Searcher.Lock.RLock()
		dir := t.Searcher.Dir()
		t.Searcher.Lock.RUnlock()

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
