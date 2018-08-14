package queue

import (
	"sync"

	"github.com/wpdirectory/wpdir/internal/metrics"
)

// Queue contains information about and helper funcs for the queue
type Queue struct {
	queue chan string
	pos   map[string]int
	sync.RWMutex
}

// New returns a Queue struct
func New(len int) *Queue {
	return &Queue{
		queue: make(chan string, len),
		pos:   make(map[string]int),
	}
}

// Add pushes an ID onto the Queue
func (q *Queue) Add(id string) {
	q.queue <- id

	q.Lock()
	q.pos[id] = len(q.queue)
	q.Unlock()

	metrics.SearchQueue.Inc()
}

// Get returns an item from the queue
// Blocks until available
func (q *Queue) Get() string {
	id := <-q.queue

	q.Lock()
	defer q.Unlock()
	delete(q.pos, id)
	q.decrementPos()

	metrics.SearchQueue.Dec()

	return id
}

// decrementPos decrements each ID by one.
func (q *Queue) decrementPos() {
	for id := range q.pos {
		q.pos[id]--
	}
}

// Pos returns the position of ID in the queue
// returns -1 if ID not found.
func (q *Queue) Pos(id string) int {
	q.RLock()
	defer q.RUnlock()
	pos, ok := q.pos[id]
	if !ok {
		return -1
	}
	return pos
}
