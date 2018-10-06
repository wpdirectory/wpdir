package tasks

import (
	"github.com/robfig/cron"
)

var runner *cron.Cron

func init() {
	runner = cron.New()
}

// Add a Task
func Add(s string, f func()) {
	runner.AddFunc(s, f)
}

// Start Task runner
func Start() {
	runner.Start()
}

// Stop Task runner
func Stop() {
	runner.Stop()
}

// Entries returns Task list
func Entries() {
	runner.Entries()
}
