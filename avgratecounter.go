package ratecounter

import (
	"strconv"
	"time"
)

// An AvgRateCounter is a thread-safe counter which returns
// the ratio between the number of calls 'Incr' and the counter value in the last interval
type AvgRateCounter struct {
	hits     *RateCounter
	counter  *RateCounter
	interval time.Duration
}

// NewRateCounter Constructs a new AvgRateCounter, for the interval provided
func NewAvgRateCounter(intrvl time.Duration) *AvgRateCounter {
	return &AvgRateCounter{
		hits:     NewRateCounter(intrvl),
		counter:  NewRateCounter(intrvl),
		interval: intrvl,
	}
}

func (r *AvgRateCounter) WithResolution(resolution int) *AvgRateCounter {
	if resolution < 1 {
		panic("AvgRateCounter resolution cannot be less than 1")
	}

	r.hits = r.hits.WithResolution(resolution)
	r.counter = r.counter.WithResolution(resolution)

	return r
}

// Incr Adds an event into the AvgRateCounter
func (a *AvgRateCounter) Incr(val int64) {
	a.hits.Incr(1)
	a.counter.Incr(val)
}

// Rate Returns the current ratio between the events count and its values during the last interval
func (a *AvgRateCounter) Rate() float64 {
	hits, value := a.hits.Rate(), a.counter.Rate()

	if hits == 0 {
		return 0 // Avoid division by zero
	}

	return float64(value) / float64(hits)
}

func (a *AvgRateCounter) String() string {
	return strconv.FormatFloat(a.Rate(), 'e', 5, 64)
}
