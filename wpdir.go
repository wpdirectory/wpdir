package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"

	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/log"
	"github.com/wpdirectory/wpdir/internal/metrics"
	"github.com/wpdirectory/wpdir/internal/server"
	"github.com/wpdirectory/wpdir/internal/tasks"
)

//go:generate go run -tags=dev embed_files.go

var (
	version string
	commit  string
	date    string
)

func main() {
	// Set and Parse flags
	flagHelp := flag.Bool("help", false, "Display help information")
	flagFresh := flag.Bool("fresh", false, "Begin with fresh data load")
	flagDev := flag.Bool("dev", false, "Enable Dev mode")
	flag.Parse()

	if *flagHelp {
		fmt.Println(helpText)
		os.Exit(1)
	}

	fmt.Println("Starting WPDirectory")

	// Create Logger
	l := log.New()

	// Create Config
	c := config.Setup(version, commit, date, *flagDev)

	l.Printf("Hostname: %s\n", c.Host)

	// Ensure Directory Structure Exists
	mkdirs(c.WD)

	// Set Temp Dir
	// TODO: Check error- what would we do?
	setTempDir(c.WD)

	// Setup Metrics
	metrics.Setup()

	// Setup BoltDB
	db.Setup(c.WD)
	defer db.Close()

	// Setup server struct to hold all App data
	s := server.New(l, c, flagFresh)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Setup HTTP server.
	s.Setup()

	// Clean up temp dir
	tasks.Add("0 */15 * * * *", emptyTempDir)

	// Start Task Runner
	tasks.Start()

	<-stop

	l.Printf("Shutting down WPdir...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Shutdown
	tasks.Stop()
	s.Shutdown(ctx)
}

const (
	helpText = `WPDirectory is a web service for lightning fast code searching of the WordPress Plugin & Theme Directories.

Usage:
  wpdir [flags]
	
Flags:
  -help      Help outputs help text and exits.
  -fresh     Begins a fresh load, all extensions are queued for updating.
  
Config:
  WPDirectory requires a config file, located at /etc/wpdir/ or in the working directory, to successfully run. See the example-config.yml.`
)

func mkdirs(wd string) {
	tmp := filepath.Join(wd, "tmp")
	os.MkdirAll(tmp, 0766)

	db := filepath.Join(wd, "data", "db")
	os.MkdirAll(db, 0766)

	ssl := filepath.Join(wd, "data", "ssl")
	os.MkdirAll(ssl, 0760)

	plugins := filepath.Join(wd, "data", "index", "plugins")
	os.MkdirAll(plugins, 0766)

	themes := filepath.Join(wd, "data", "index", "themes")
	os.MkdirAll(themes, 0766)
}

// setTempDir sets the temp dir
// WPdir creates a lot of temp files which need
// to be cleaned up by the program itself
func setTempDir(wd string) error {
	path := filepath.Join(wd, "tmp")

	switch opsys := runtime.GOOS; opsys {
	case "windows":
		err := os.Setenv("TMP", path)
		return err
	case "darwin":
		err := os.Setenv("TMPDIR", path)
		return err
	case "linux":
		err := os.Setenv("TMPDIR", path)
		return err
	default:
		return nil
	}
}

// emptyTempDir deletes the contents of the temp dir
// Clears out temp files created during indexing
func emptyTempDir() {
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "tmp")

	d, err := os.Open(dir)
	if err != nil {
		return
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return
		}
	}

	return
}
