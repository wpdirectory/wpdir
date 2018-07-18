// +build dev

package data

//go:generate go run -tags=dev assets_generate.go

import (
	"net/http"
	"path/filepath"
)

// Assets contains project assets.
var Assets http.FileSystem = http.Dir(filepath.Join("web", "build"))
