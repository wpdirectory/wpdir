package limit

import (
	"time"

	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/store/memory"
)

var instance *limiter.Limiter

func init() {
	rate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  20,
	}
	store := memory.NewStore()

	instance = limiter.New(store, rate)
}

// New returns a new limiter
func New() *limiter.Limiter {
	return instance
}
