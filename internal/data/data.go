// +build dev

package data

//go:generate go run -tags=dev assets_generate.go

import (
	"net/http"
	"path/filepath"
)

// Assets contains project assets.
var Assets http.FileSystem = http.Dir(filepath.Join("D:/projects/go/src/github.com/wpdirectory/wpdir", "web", "build"))
