package ratecounter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAvgRateCounter(t *testing.T) {
	interval := 50 * time.Millisecond
	r := NewAvgRateCounter(interval)

	assert.Equal(t, float64(0), r.Rate())
	assert.Equal(t, int64(0), r.Hits())
	r.Incr(1) // counter = 1, hits = 1
	assert.Equal(t, float64(1.0), r.Rate())
	assert.Equal(t, int64(1), r.Hits())
	r.Incr(3) // counter = 4, hits = 2
	assert.Equal(t, float64(2.0), r.Rate())
	assert.Equal(t, int64(2), r.Hits())
	time.Sleep(2 * interval)
	assert.Equal(t, float64(0), r.Rate())
	assert.Equal(t, int64(0), r.Hits())
}

func TestAvgRateCounterAdvanced(t *testing.T) {
	interval := 50 * time.Millisecond
	almost := 45 * time.Millisecond
	gap := 1 * time.Millisecond
	r := NewAvgRateCounter(interval)

	assert.Equal(t, float64(0), r.Rate())
	assert.Equal(t, int64(0), r.Hits())
	r.Incr(1) // counter = 1, hits = 1
	assert.Equal(t, float64(1.0), r.Rate())
	assert.Equal(t, int64(1), r.Hits())
	time.Sleep(interval - almost)
	r.Incr(3) // counter = 4, hits = 2
	assert.Equal(t, float64(2.0), r.Rate())
	assert.Equal(t, int64(2), r.Hits())
	time.Sleep(almost + gap)
	assert.Equal(t, float64(3.0), r.Rate())
	assert.Equal(t, int64(1), r.Hits()) // counter = 3, hits = 1
	time.Sleep(2 * interval)
	assert.Equal(t, float64(0), r.Rate())
	assert.Equal(t, int64(0), r.Hits())
}

func TestAvgRateCounterMinResolution(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Resolution < 1 did not panic")
		}
	}()

	NewAvgRateCounter(500 * time.Millisecond).WithResolution(0)
}

func TestAvgRateCounterNoResolution(t *testing.T) {
	interval := 50 * time.Millisecond
	almost := 45 * time.Millisecond
	gap := 1 * time.Millisecond
	r := NewAvgRateCounter(interval).WithResolution(1)

	assert.Equal(t, float64(0), r.Rate())
	assert.Equal(t, int64(0), r.Hits())
	r.Incr(1) // counter = 1, hits = 1
	assert.Equal(t, float64(1.0), r.Rate())
	assert.Equal(t, int64(1), r.Hits())
	time.Sleep(interval - almost)
	r.Incr(3) // counter = 4, hits = 2
	assert.Equal(t, float64(2.0), r.Rate())
	assert.Equal(t, int64(2), r.Hits())
	time.Sleep(almost + gap)
	assert.Equal(t, float64(0), r.Rate())
	assert.Equal(t, int64(0), r.Hits()) // counter = 0, hits = 0, r.Hits())
	time.Sleep(2 * interval)
	assert.Equal(t, float64(0), r.Rate())
	assert.Equal(t, int64(0), r.Hits())
}

func TestAvgRateCounter_String(t *testing.T) {
	r := NewAvgRateCounter(1 * time.Second)
	if r.String() != "0.00000e+00" {
		t.Error("Expected ", r.String(), " to equal ", "0.00000e+00")
	}

	r.Incr(1)
	if r.String() != "1.00000e+00" {
		t.Error("Expected ", r.String(), " to equal ", "1.00000e+00")
	}
}

func TestAvgRateCounter_Incr_ReturnsImmediately(t *testing.T) {
	interval := 1 * time.Second
	r := NewAvgRateCounter(interval)

	start := time.Now()
	r.Incr(-1)
	duration := time.Since(start)

	if duration >= 1*time.Second {
		t.Error("incr took", duration, "to return")
	}
}

func BenchmarkAvgRateCounter(b *testing.B) {
	interval := 1 * time.Millisecond
	r := NewAvgRateCounter(interval)

	for i := 0; i < b.N; i++ {
		r.Incr(1)
		r.Rate()
	}
}
