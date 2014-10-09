package ratecounter

import (
	"sync/atomic"
	"time"
)

// A Counter is a thread-safe counter implementation
type Counter struct {
	value *int64
}

// NewCounter is used to construct a new Counter object
func NewCounter() *Counter {
	var v int64
	return &Counter{
		value: &v,
	}
}

// Increment the counter by some value
func (c *Counter) Incr(val int64) {
	atomic.AddInt64(c.value, val)
}

// Return the counter's current value
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(c.value)
}

// A RateCounter is a thread-safe counter which returns the number of times
// 'Incr' has been called in the last interval
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

// Add an event into the RateCounter
func (r *RateCounter) Incr(val int64) {
	r.counter.Incr(val)
	r.scheduleDecrement(-1 * val)
}

func (r *RateCounter) scheduleDecrement(amount int64) {
	time.AfterFunc(r.interval, func() {
		r.counter.Incr(amount)
	})
}

// Return the current number of events in the last interval
func (r *RateCounter) Rate() int64 {
	return r.counter.Value()
}
