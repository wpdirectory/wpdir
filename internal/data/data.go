// +build dev

package data

import (
	"net/http"
	"path/filepath"
)

// Assets contains project assets.
var Assets http.FileSystem = http.Dir(filepath.Join("../../", "web", "build"))
