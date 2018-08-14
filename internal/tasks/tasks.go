package tasks

import (
	"github.com/robfig/cron"
)

var runner *cron.Cron

func init() {
	runner = cron.New()
}

// Add ...
func Add(s string, f func()) {
	runner.AddFunc(s, f)
}

// Start ...
func Start() {
	runner.Start()
}

// Stop ...
func Stop() {
	runner.Stop()
}

// Entries ...
func Entries() {
	runner.Entries()
}
