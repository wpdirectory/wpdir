package store

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/svn"
)

// DeletePlugin deletes the plugin fromstorage.
// Moved from using os.RemoveAll() due to weaknesses with high depth
func DeletePlugin(slug string) error {

	path := filepath.Join(storageDir, "plugins", slug)
	_, err := exec.Command("rm", "-rf", path).Output()

	return err

}

// AddPlugin adds the plugin to storage.
func AddPlugin(remotePath string, localPath string) error {

	dest := filepath.Join(storageDir, "plugins", localPath)

	err := repo.DoExport("plugins", remotePath, dest)

	return err

}

// UpdatePlugin updates the plugin in storage.
func UpdatePlugin(remotePath string, localPath string) error {

	err := DeletePlugin(localPath)
	if err != nil {
		return err
	}

	err = AddPlugin(remotePath, localPath)
	if err != nil {
		return err
	}

	return nil

}

// ListPlugins returns a list of all plugins in storage.
func ListPlugins() ([]string, error) {

	path := filepath.Join(storageDir, "plugins")

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var plugins []string
	for _, file := range files {
		plugins = append(plugins, file.Name())
	}

	return plugins, nil

}

// FreshStart ...
func FreshStart() error {

	list, err := repo.GetList("plugins", "")
	if err != nil {
		return err
	}

	limiter := make(chan struct{}, 5)

	list = []svn.ListEntry{
		svn.ListEntry{
			Name: "kebo-twitter-feed",
		},
		svn.ListEntry{
			Name: "kebo-social",
		},
		svn.ListEntry{
			Name: "health-check",
		},
		svn.ListEntry{
			Name: "wordpress-seo",
		},
		svn.ListEntry{
			Name: "wordfence",
		},
		svn.ListEntry{
			Name: "akismet",
		},
	}

	var wg sync.WaitGroup

	for _, item := range list {

		// Will block if more than max Goroutines already running.
		limiter <- struct{}{}
		wg.Add(1)

		go func(name string, wg *sync.WaitGroup) {

			log.Printf("%s (plugin) is being updated.\n", name)

			// Update Plugin Files
			err := UpdatePlugin(name+"/trunk/", name)
			if err != nil {
				log.Printf("%s (plugin) could not be updated: %s\n", name, err)
			}

			log.Printf("%s (plugin) is being indexed.\n", name)

			// Update Plugin Index
			opts := &index.IndexOptions{
				ExcludeDotFiles: true,
			}
			wd, _ := os.Getwd()
			src := filepath.Join(wd, "data", "plugins", name)
			dest := filepath.Join(wd, "index", "plugins", name)
			url := name
			rev := "3462396423894"
			_, err = index.BuildNew(opts, dest, src, url, string(rev))
			if err != nil {
				log.Printf("%s (plugin) could not be indexed: %s\n", name, err)
			}

			<-limiter
			wg.Done()

		}(item.Name, &wg)

	}

	wg.Wait()

	return nil

}
