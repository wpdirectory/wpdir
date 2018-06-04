package sync

import (
	"log"
	"os"

	"github.com/wpdirectory/wpdir/internal/svn"
)

var (
	storageDir  string
	indexDir    string
	workerCount int
	bufferSize  int
	wd          string
)

// Setup ...
func Setup() {

	// Check if SVN client is available
	if !svn.IsClientInstalled() {
		log.Fatal("The SVN client is not available, it is required for operation.")
	}

	// Setup Config Based Values
	storageDir = "data"
	indexDir = "index"
	workerCount = 10
	bufferSize = 50
	wd, _ = os.Getwd()

	// Create Plugins Channel
	plugins := make(chan svn.LogEntry, bufferSize)

	// One Gorountine to monitor Plugins changelog
	go monitorPlugins(plugins)

	// Many Gorountines checking if Plugins need testing
	go startPluginChecks(plugins)

	// Create Themes Channel
	//themes := make(chan string, bufferSize)

	// One Gorountine to monitor Themes changelog
	//go monitorThemes(themes)

	// Many Gorountines checking if Themes need testing
	//go startThemeChecks(themes)

}
