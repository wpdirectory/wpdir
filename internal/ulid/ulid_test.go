package ulid

import (
	"testing"

	"github.com/oklog/ulid"
)

func TestNew(t *testing.T) {
	var ulids []string
	for i := 0; i < 3; i++ {
		ulids = append(ulids, New())
	}

	for _, id := range ulids {
		_, err := ulid.Parse(id)
		if err != nil {
			t.Errorf("Not a valid ULID: %s\n", err)
		}
	}
}
