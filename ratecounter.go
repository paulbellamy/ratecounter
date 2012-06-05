// Copyright 2012 Paul Bellamy. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ratecounter

import (
	"time"
)

// A Counter is a thread-safe counter implementation
type Counter struct {
	value int64
	lock  chan int
}

// NewCounter is used to construct a new Counter object
func NewCounter() *Counter {
	return &Counter{
		value: 0,
		lock:  make(chan int, 1),
	}
}

// Increment the counter by some value
func (c *Counter) Incr(val int64) {
	c.lock <- 1
	c.value += val
	<-c.lock
}

// Return the counter's current value
func (c *Counter) Value() int64 {
	c.lock <- 1
	val := c.value
	<-c.lock
	return val
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
