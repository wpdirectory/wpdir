package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/log"
	"github.com/wpdirectory/wpdir/internal/server"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	// Set and Parse flags
	flagHelp := flag.Bool("help", false, "Display help information")
	flagFresh := flag.Bool("fresh", false, "Begin with fresh data load")
	flag.Parse()

	if *flagHelp {
		fmt.Println(helpText)
		os.Exit(1)
	}

	fmt.Println("Starting WPDirectory")

	// Create Logger
	l := log.New()

	// Create Config
	c := config.Setup(version, commit, date)

	l.Printf("Hostname: %s\n", c.Host)

	//l.Printf("Version: %s Commit: %s Date: %s\n", Version, Commit, Date)

	// Ensure Directory Structure Exists
	mkdirs(c.WD)

	// Setup BoltDB
	db.Setup(c.WD)
	defer db.Close()

	// Setup server struct to hold all App data
	s := server.New(l, c, flagFresh)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Setup HTTP server.
	s.Setup()

	<-stop

	l.Printf("Shutting down WPdir...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.Shutdown(ctx)
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

func mkdirs(wd string) {
	db := filepath.Join(wd, "data", "db")
	os.MkdirAll(db, os.ModePerm)

	plugins := filepath.Join(wd, "data", "index", "plugins")
	os.MkdirAll(plugins, os.ModePerm)

	themes := filepath.Join(wd, "data", "index", "themes")
	os.MkdirAll(themes, os.ModePerm)
}
