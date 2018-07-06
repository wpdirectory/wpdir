package files

import (
	"archive/zip"
	"path/filepath"
	"strings"
	"sync"
)

// Stats contains
type Stats struct {
	Files      []File
	TotalFiles int
	TotalSize  int64
	Summary    Summary
	sync.RWMutex
}

// File contains basic data about a specific file
type File struct {
	Name      string
	Extension string
	Size      int64
}

// Summary contains an overview of the PHP/JS and CSS files contained
type Summary struct {
	PHP uint8
	JS  uint8
	CSS uint8
}

// New returns an empty FileStats struct
func New() *Stats {
	return &Stats{}
}

// AddFile adds a file to the files field
// Only stores Name, Extension and Size
func (s *Stats) AddFile(zf *zip.File) {
	f := zf.FileInfo()
	if f.IsDir() {
		return
	}

	s.RLock()
	defer s.RUnlock()
	file := File{
		Name:      f.Name(),
		Extension: strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
		Size:      f.Size(),
	}
	s.TotalFiles++
	s.TotalSize += f.Size()
	s.Files = append(s.Files, file)
}

// GenerateSummary creates a Summary using data from the Files field
func (s *Stats) GenerateSummary() {
	if len(s.Files) == 0 {
		return
	}
	s.RLock()
	defer s.RUnlock()
	var php, js, css, total int64
	for _, file := range s.Files {
		switch file.Extension {
		case "php":
			php += file.Size
			total += file.Size
			break
		case "js":
			js += file.Size
			total += file.Size
			break
		case "css":
			css += file.Size
			total += file.Size
			break
		}
	}
	if total == 0 {
		return
	}
	s.Summary = Summary{
		PHP: uint8((php / total) * 100),
		JS:  uint8((js / total) * 100),
		CSS: uint8((css / total) * 100),
	}
}
