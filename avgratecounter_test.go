package ratecounter

import (
	"testing"
	"time"
)

func TestAvgRateCounter(t *testing.T) {
	interval := 500 * time.Millisecond
	r := NewAvgRateCounter(interval)

	check := func(expected float64) {
		val := r.Rate()
		if val != expected {
			t.Error("Expected ", val, " to equal ", expected)
		}
	}

	check(0)
	r.Incr(1) // counter = 1, hits = 1
	check(1.0)
	r.Incr(3) // counter = 4, hits = 2
	check(2.0)
	time.Sleep(2 * interval)
	check(0)
}

func TestAvgRateCounterAdvanced(t *testing.T) {
	interval := 500 * time.Millisecond
	almost := 450 * time.Millisecond
	r := NewAvgRateCounter(interval)

	check := func(expected float64) {
		val := r.Rate()
		if val != expected {
			t.Error("Expected ", val, " to equal ", expected)
		}
	}

	check(0)
	r.Incr(1) // counter = 1, hits = 1
	check(1.0)
	time.Sleep(interval - almost)
	r.Incr(3) // counter = 4, hits = 2
	check(2.0)
	time.Sleep(almost)
	check(3.0) // counter = 3, hits = 1
	time.Sleep(2 * interval)
	check(0)
}

func TestAvgRateCounterNoResolution(t *testing.T) {
	interval := 500 * time.Millisecond
	almost := 450 * time.Millisecond
	r := NewAvgRateCounter(interval).WithResolution(1)

	check := func(expected float64) {
		val := r.Rate()
		if val != expected {
			t.Error("Expected ", val, " to equal ", expected)
		}
	}

	check(0)
	r.Incr(1) // counter = 1, hits = 1
	check(1.0)
	time.Sleep(interval - almost)
	r.Incr(3) // counter = 4, hits = 2
	check(2.0)
	time.Sleep(almost)
	check(0) // counter = 0, hits = 0
	time.Sleep(2 * interval)
	check(0)
}

func TestAvgRateCounter_Incr_ReturnsImmediately(t *testing.T) {
	interval := 1 * time.Second
	r := NewRateCounter(interval)

	start := time.Now()
	r.Incr(-1)
	duration := time.Since(start)

	if duration >= 1*time.Second {
		t.Error("incr took", duration, "to return")
	}
}

func BenchmarkAvgRateCounter(b *testing.B) {
	interval := 0 * time.Millisecond
	r := NewAvgRateCounter(interval)

	for i := 0; i < b.N; i++ {
		r.Incr(1)
		r.Rate()
	}
}
