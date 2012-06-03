// Package rate_counter provides a concurrent rate-counter, for tracking counts
// in an interval
package rate_counter

import (
	"time"
  "sync"
)

// A Counter is a thread-safe counter implementation
type Counter struct {
	value int64
  mutex *sync.Mutex
}

// NewCounter is used to construct a new Counter object
func NewCounter() *Counter {
	return &Counter{
		value: 0,
    mutex: &sync.Mutex{},
	}
}

// Increment the counter by some value
func (c *Counter) Incr(val int64) {
  c.mutex.Lock()
  defer c.mutex.Unlock()
  c.value += val
}

// Return the counter's current value
func (c *Counter) Value() int64 {
  c.mutex.Lock()
  defer c.mutex.Unlock()
  return c.value
}

// A RateCounter is a thread-safe counter which returns the number of times
// 'Mark' has been called in the last interval
type RateCounter struct {
	counter  *Counter
	interval time.Duration
}

// Constructs a new RateCounter, for the interval provided
func NewRateCounter(intrvl time.Duration) *RateCounter {
	return &RateCounter{
		counter:  NewCounter(),
		interval: intrvl,
	}
}

// Add 1 event into the RateCounter
func (r *RateCounter) Mark() {
	r.counter.Incr(1)
	time.AfterFunc(r.interval, func() {
		r.counter.Incr(-1)
	})
}

// Return the current number of events in the last interval
func (r *RateCounter) Rate() int64 {
	return r.counter.Value()
}
