package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/log"
	"github.com/wpdirectory/wpdir/internal/server"
	"github.com/wpdirectory/wpdir/internal/stats"
	"github.com/wpdirectory/wpdir/internal/store"
	"github.com/wpdirectory/wpdir/internal/sync"
)

func main() {

	// Set and Parse flags
	flagHelp := flag.Bool("help", false, "")
	flag.Parse()

	if *flagHelp {
		fmt.Println(helpText)
		os.Exit(1)
	}

	// Setup Stats.
	stats.Setup()

	// Setup server struct to hold all App data
	s := server.New()

	// Setup Logger
	s.Logger = log.New()

	// Setup Config
	s.Config = config.Setup()

	// Setup Data Store
	s.Store = store.New(s.Config)

	s.Logger.Printf("Starting WPDirectory - Version: %s\n", s.Config.Version)

	// Setup background sync process.
	sync.Setup()

	//err := store.FreshStart()
	//if err != nil {
	//s.Logger.Fatalf("Failed fresh start: %s\n", err)
	//}

	// Setup HTTP server.
	s.Setup()

}

const (
	helpText = `WPDirectory is a web service for lightning fast code searching of the WordPress Plugin & Theme Directories.

Usage:
  wpdir [flags]

Version: 0.5.0
	
Flags:
  --help      Help outputs help text and exits.
  
Config:
  WPDirectory requires a config file, located at /etc/wpdir/ or in the working directory, to successfully run. See the example-config.yml.`
)
