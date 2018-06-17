package store

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/wpdirectory/wpdir/internal/repo"
)

// DeletePlugin deletes the plugin fromstorage.
// Moved from using os.RemoveAll() due to weaknesses with high depth
func DeletePlugin(slug string) error {

	path := filepath.Join(storageDir, "tmp", "plugins", slug)
	_, err := exec.Command("rm", "-rf", path).Output()

	return err

}

// AddPlugin adds the plugin to storage.
func AddPlugin(remotePath string, localPath string) error {

	dest := filepath.Join(storageDir, "tmp", "plugins", localPath)

	err := repo.DoExport("plugins", remotePath, dest)

	return err

}

// GetPlugin fetches the plugin files to local temp directory.
func GetPlugin(remotePath string, localPath string) error {

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

	path := filepath.Join(storageDir, "tmp", "plugins")

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
