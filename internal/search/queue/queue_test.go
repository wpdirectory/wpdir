package queue

import (
	"testing"

	"github.com/wpdirectory/wpdir/internal/metrics"
)

func init() {
	metrics.Setup()
}

func TestAdd(t *testing.T) {
	queue := &Queue{
		queue: make(chan string, 10),
		pos:   make(map[string]int),
	}
	ids := []string{
		"01CH6CNP575QSN84B4YH1FRGYC",
		"01CH5K5GW8BQJ1CSTM2B4EZ5SZ",
		"01CH40V4G0DZ2GRNAQ3Z3FT85H",
	}

	for key, id := range ids {
		queue.Add(id)
		got := len(queue.queue)
		want := key + 1
		if got != want {
			t.Errorf("Expected length %d got %d", want, got)
		}
	}
}

func TestGet(t *testing.T) {
	queue := &Queue{
		queue: make(chan string, 10),
		pos:   make(map[string]int),
	}
	ids := []string{
		"01CH6CNP575QSN84B4YH1FRGYC",
		"01CH5K5GW8BQJ1CSTM2B4EZ5SZ",
		"01CH40V4G0DZ2GRNAQ3Z3FT85H",
	}

	for _, id := range ids {
		queue.Add(id)
	}

	for _, id := range ids {
		got := queue.Get()
		want := id
		if got != want {
			t.Errorf("Expected ID %s got %s", want, got)
		}
	}
}

func TestPos(t *testing.T) {
	queue := &Queue{
		queue: make(chan string, 10),
		pos:   make(map[string]int),
	}
	ids := []string{
		"01CH6CNP575QSN84B4YH1FRGYC",
		"01CH5K5GW8BQJ1CSTM2B4EZ5SZ",
		"01CH40V4G0DZ2GRNAQ3Z3FT85H",
	}

	for _, id := range ids {
		queue.Add(id)
	}

	for key, id := range ids {
		got := queue.Pos(id)
		want := key + 1
		if got != want {
			t.Errorf("Expected position %d got %d", want, got)
		}
	}
}

func TestDecrementPos(t *testing.T) {
	queue := &Queue{
		queue: make(chan string, 10),
		pos:   make(map[string]int),
	}

	ids := []string{
		"01CH6CNP575QSN84B4YH1FRGYC",
		"01CH5K5GW8BQJ1CSTM2B4EZ5SZ",
		"01CH40V4G0DZ2GRNAQ3Z3FT85H",
	}

	for _, id := range ids {
		queue.Add(id)
	}

	id := queue.Get()
	got := queue.Pos(id)
	want := -1
	if got != want {
		t.Errorf("Expected invalid key (%d) got %d", want, got)
	}

	ids = ids[1:]

	for key, id := range ids {
		got := queue.Pos(id)
		want := key + 1
		if got != want {
			t.Errorf("Expected position %d got %d", want, got)
		}
	}
}
