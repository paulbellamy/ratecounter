package ratecounter

import (
	"strconv"
	"sync/atomic"
	"time"
)

// An AvgRateCounter is a thread-safe counter which returns
// the ratio between the number of calls 'Incr' and the counter value in the last interval
type AvgRateCounter struct {
	hits     int64
	counter  Counter
	interval time.Duration
}

// NewRateCounter Constructs a new AvgRateCounter, for the interval provided
func NewAvgRateCounter(intrvl time.Duration) *AvgRateCounter {
	return &AvgRateCounter{
		interval: intrvl,
	}
}

// Incr Adds an event into the AvgRateCounter
func (a *AvgRateCounter) Incr(val int64) {
	atomic.AddInt64(&a.hits, 1)
	a.counter.Incr(val)

	time.AfterFunc(a.interval, func() {
		atomic.AddInt64(&a.hits, -1)
		a.counter.Incr(-1 * val)
	})
}

// Rate Returns the current ratio between the events count and its values during the last interval
func (a *AvgRateCounter) Rate() float64 {
	hits, value := atomic.LoadInt64(&a.hits), a.counter.Value()

	if hits == 0 {
		return 0 // Avoid division by zero
	}

	return float64(value) / float64(hits)
}

func (a *AvgRateCounter) String() string {
	return strconv.FormatFloat(a.Rate(), 'e', 5, 64)
}
