package store

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/wpdirectory/wpdir/internal/config"
)

var storageDir string

// DataStore ...
type DataStore interface {
	CacheSet(string, interface{}, time.Duration)
	CacheGet(string) (interface{}, bool)
	CacheDelete(string)

	Close()
}

// Store ...
type Store struct {
	cache *cache.Cache
}

// New ...
func New(c *config.Config) DataStore {

	dataDir := c.DataDir

	if !path.IsAbs(dataDir) {

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not establish working directory, failed to setup storage.")
		}

		dataDir := filepath.Join(wd, dataDir)

		err = os.MkdirAll(dataDir, os.ModeDir)
		if err != nil && err != os.ErrExist {
			log.Fatal("Could not create storage directory, failed to setup storage.")
		}

	}

	storageDir = dataDir

	return &Store{
		cache: cache.New(cache.NoExpiration, 24*time.Hour),
	}
}

// Close ends the database connections.
func (s *Store) Close() {

	//s.sql.Close()

}

// DeleteFolder deletes the folder and all its contents at path.
func DeleteFolder(filepath string) error {
	_, err := exec.Command("rm", "-rf", filepath).Output()

	return err
}
