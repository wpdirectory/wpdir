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

func main() {

	// Set and Parse flags
	flagHelp := flag.Bool("help", false, "")
	flag.Parse()

	if *flagHelp {
		fmt.Println(helpText)
		os.Exit(1)
	}

	fmt.Println("Starting WPDirectory")

	// Setup Stats.
	//stats.Setup()

	// Create Logger
	l := log.New()

	// Create Config
	c := config.Setup()

	// Ensure Directory Structure Exists
	mkdir(c.WD)

	//err := certs.Get(c.DNS.Email, c.DNS.APIKey)
	//panic(err)

	// Setup BoltDB
	db.Setup(c.WD)
	defer db.Close()

	// Setup server struct to hold all App data
	s := server.New(l, c)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Setup HTTP server.
	s.Setup()

	<-stop

	l.Printf("Shutting down WPdir...\n")

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
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

func mkdir(wd string) {

	db := filepath.Join(wd, "data", "db")
	os.MkdirAll(db, os.ModePerm)

	plugins := filepath.Join(wd, "data", "index", "plugins")
	os.MkdirAll(plugins, os.ModePerm)

	themes := filepath.Join(wd, "data", "index", "themes")
	os.MkdirAll(themes, os.ModePerm)

}
