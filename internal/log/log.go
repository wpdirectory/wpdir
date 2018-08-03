package log

import (
	"log"
	"os"
)

// New returns a new logger
func New() *log.Logger {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime)
}
