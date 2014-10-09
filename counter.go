package ratecounter

import "sync/atomic"

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
