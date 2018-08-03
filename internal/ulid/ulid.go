package ulid

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

var entropy *rand.Rand
var mutex sync.RWMutex

func init() {
	entropy = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// New returns a new ULID
func New() string {
	mutex.RLock()
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	mutex.RUnlock()

	return id.String()
}
