package ratecounter

import (
	"sync/atomic"
	"time"
)

// A RateCounter is a thread-safe counter which returns the number of times
// 'Incr' has been called in the last interval
type ResettableRateCounter struct {
	ptr atomic.Pointer[RateCounter]
}

// NewResettableRateCounter Constructs a new ResettableRateCounter, for the interval provided
func NewResettableRateCounter(interval time.Duration) *ResettableRateCounter {
	return Resettable(NewRateCounter(interval))
}

// NewResettableRateCounterWithResolution Constructs a new ResettableRateCounter, for the provided interval and resolution
func NewResettableRateCounterWithResolution(interval time.Duration, resolution int) *ResettableRateCounter {
	return Resettable(NewRateCounterWithResolution(interval, resolution))
}

// Resettable wraps a RateCounter to allow resetting it. Note, this adds a small overhead to each operation.
func Resettable(rc *RateCounter) *ResettableRateCounter {
	r := &ResettableRateCounter{}
	r.ptr.Store(rc)
	return r
}

// WithResolution determines the minimum resolution of this counter, default is 20
func (r *ResettableRateCounter) WithResolution(resolution int) *ResettableRateCounter {
	r.ptr.Load().WithResolution(resolution)
	return r
}

// OnStop allow to specify a function that will be called each time the counter
// reaches 0. Useful for removing it.
func (r *ResettableRateCounter) OnStop(f func(*RateCounter)) {
	r.ptr.Load().OnStop(f)
}

// Incr Add an event into the RateCounter
func (r *ResettableRateCounter) Incr(val int64) {
	r.ptr.Load().Incr(val)
}

// Rate Return the current number of events in the last interval
func (r *ResettableRateCounter) Rate() int64 {
	return r.ptr.Load().Rate()
}

// MaxRate counts the maximum instantaneous change in rate.
//
// This is useful to calculate number of events in last period without
// "averaging" effect. i.e. currently if counter is set for 30 seconds
// duration, and events fire 10 times per second, it'll take 30 seconds for
// "Rate" to show 300 (or 10 per second). The "MaxRate" will show 10
// immediately, and it'll stay this way for the next 30 seconds, even if rate
// drops below it.
func (r *ResettableRateCounter) MaxRate() int64 {
	return r.ptr.Load().MaxRate()
}

func (r *ResettableRateCounter) String() string {
	return r.ptr.Load().String()
}

// Reset method resets the counter's values to zero, like it was just
// initialized.
func (r *ResettableRateCounter) Reset() {
	old := r.ptr.Load()
	r.ptr.Store(NewRateCounterWithResolution(old.interval, old.resolution))
}
