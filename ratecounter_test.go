package ratecounter

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateCounter(t *testing.T) {
	interval := 50 * time.Millisecond
	r := NewRateCounter(interval)

	assert.Equal(t, int64(0), r.Rate())
	r.Incr(1)
	assert.Equal(t, int64(1), r.Rate())
	r.Incr(2)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(2 * interval)
	assert.Equal(t, int64(0), r.Rate())
}

func TestRateCounterExpireAndRestart(t *testing.T) {
	interval := 50 * time.Millisecond

	r := NewRateCounter(interval)

	assert.Equal(t, int64(0), r.Rate())
	r.Incr(1)
	assert.Equal(t, int64(1), r.Rate())

	// Let it expire down to zero, then restart
	time.Sleep(2 * interval)
	assert.Equal(t, int64(0), r.Rate())
	time.Sleep(2 * interval)
	r.Incr(2)
	assert.Equal(t, int64(2), r.Rate())

	// Let it expire down to zero
	time.Sleep(2 * interval)
	assert.Equal(t, int64(0), r.Rate())

	// Restart it
	r.Incr(2)
	assert.Equal(t, int64(2), r.Rate())
}

func TestRateCounterPartial(t *testing.T) {
	interval := 50 * time.Millisecond
	almost := 40 * time.Millisecond

	r := NewRateCounter(interval)

	assert.Equal(t, int64(0), r.Rate())
	r.Incr(1)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(almost)
	r.Incr(2)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(almost)
	assert.Equal(t, int64(2), r.Rate())
	time.Sleep(2 * interval)
	assert.Equal(t, int64(0), r.Rate())
}

func TestRateCounterHighResolution(t *testing.T) {
	interval := 50 * time.Millisecond
	tenth := 5 * time.Millisecond

	r := NewRateCounter(interval).WithResolution(100)

	assert.Equal(t, int64(0), r.Rate())
	r.Incr(1)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(2 * tenth)
	r.Incr(1)
	assert.Equal(t, int64(2), r.Rate())
	time.Sleep(2 * tenth)
	r.Incr(1)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(interval - 5*tenth)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(2), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.Rate())
}

func TestRateCounterLowResolution(t *testing.T) {
	interval := 50 * time.Millisecond
	tenth := 5 * time.Millisecond

	r := NewRateCounter(interval).WithResolution(4)

	assert.Equal(t, int64(0), r.Rate())
	r.Incr(1)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(2 * tenth)
	r.Incr(1)
	assert.Equal(t, int64(2), r.Rate())
	time.Sleep(2 * tenth)
	r.Incr(1)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(interval - 5*tenth)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.Rate())
}

func TestNewRateCounterWithResolution(t *testing.T) {
	interval := 50 * time.Millisecond
	tenth := 5 * time.Millisecond

	r := NewRateCounterWithResolution(interval, 4)

	// Same as previous test with low resolution
	assert.Equal(t, int64(0), r.Rate())
	r.Incr(1)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(2 * tenth)
	r.Incr(1)
	assert.Equal(t, int64(2), r.Rate())
	time.Sleep(2 * tenth)
	r.Incr(1)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(interval - 5*tenth)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.Rate())
}

func TestRateCounterMinResolution(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Resolution < 1 did not panic")
		}
	}()

	NewRateCounter(50 * time.Millisecond).WithResolution(0)
}

func TestRateCounterNoResolution(t *testing.T) {
	interval := 50 * time.Millisecond
	tenth := 5 * time.Millisecond

	r := NewRateCounter(interval).WithResolution(1)

	assert.Equal(t, int64(0), r.Rate())
	r.Incr(1)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(2 * tenth)
	r.Incr(1)
	assert.Equal(t, int64(2), r.Rate())
	time.Sleep(2 * tenth)
	r.Incr(1)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(interval - 5*tenth)
	assert.Equal(t, int64(3), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.Rate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.Rate())
}

func TestRateCounter_String(t *testing.T) {
	r := NewRateCounter(1 * time.Second)
	if r.String() != "0" {
		t.Error("Expected ", r.String(), " to equal ", "0")
	}

	r.Incr(1)
	if r.String() != "1" {
		t.Error("Expected ", r.String(), " to equal ", "1")
	}
}

func TestRateCounterHighResolutionMaxRate(t *testing.T) {
	interval := 500 * time.Millisecond
	tenth := 50 * time.Millisecond

	r := NewRateCounter(interval).WithResolution(100)

	assert.Equal(t, int64(0), r.MaxRate())
	r.Incr(3)
	assert.Equal(t, int64(3), r.MaxRate())
	time.Sleep(2 * tenth)
	r.Incr(2)
	assert.Equal(t, int64(3), r.MaxRate())
	time.Sleep(2 * tenth)
	r.Incr(4)
	assert.Equal(t, int64(4), r.MaxRate())
	time.Sleep(interval - 5*tenth)
	assert.Equal(t, int64(4), r.MaxRate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(4), r.MaxRate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(4), r.MaxRate())
	time.Sleep(2 * tenth)
	assert.Equal(t, int64(0), r.MaxRate())
}

func TestRateCounter_Incr_ReturnsImmediately(t *testing.T) {
	interval := 1 * time.Second
	r := NewRateCounter(interval)

	start := time.Now()
	r.Incr(-1)
	duration := time.Since(start)

	if duration >= 1*time.Second {
		t.Error("incr took", duration, "to return")
	}
}

func TestRateCounter_OnStop(t *testing.T) {
	var called Counter
	interval := 50 * time.Millisecond
	r := NewRateCounter(interval)
	r.OnStop(func(r *RateCounter) {
		called.Incr(1)
	})
	r.Incr(1)

	current := called.Value()
	if current != 0 {
		t.Error("Expected called to equal 0, got ", current)
	}

	time.Sleep(2 * interval)
	current = called.Value()
	if current != 1 {
		t.Error("Expected called to equal 1, got ", current)
	}
}

func BenchmarkRateCounter(b *testing.B) {
	interval := 1 * time.Millisecond
	r := NewRateCounter(interval)

	for i := 0; i < b.N; i++ {
		r.Incr(1)
		r.Rate()
	}
}

func BenchmarkRateCounter_Parallel(b *testing.B) {
	interval := 1 * time.Millisecond
	r := NewRateCounter(interval)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Incr(1)
			r.Rate()
		}
	})
}

func BenchmarkRateCounter_With5MillionExisting(b *testing.B) {
	interval := 1 * time.Hour
	r := NewRateCounter(interval)

	for i := 0; i < 5000000; i++ {
		r.Incr(1)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Incr(1)
		r.Rate()
	}
}

func Benchmark_TimeNowAndAdd(b *testing.B) {
	var a time.Time
	for i := 0; i < b.N; i++ {
		a = time.Now().Add(1 * time.Second)
	}
	fmt.Fprintln(ioutil.Discard, a)
}
