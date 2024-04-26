package ratecounter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResettableRateCounterResetAndRestart(t *testing.T) {
	interval := 50 * time.Millisecond

	r := NewResettableRateCounter(interval)

	assert.Equal(t, int64(0), r.Rate())
	r.Incr(1)
	assert.Equal(t, int64(1), r.Rate())
	time.Sleep(2 * interval)
	assert.Equal(t, int64(0), r.Rate())
	time.Sleep(2 * interval)
	r.Incr(2)
	assert.Equal(t, int64(2), r.Rate())
	time.Sleep(2 * interval)
	assert.Equal(t, int64(0), r.Rate())
	r.Incr(2)
	assert.Equal(t, int64(2), r.Rate())
	r.Reset()
	assert.Equal(t, int64(0), r.Rate())
	r.Incr(3)
	assert.Equal(t, int64(3), r.Rate())
}

func BenchmarkResettableRateCounter(b *testing.B) {
	interval := 1 * time.Millisecond
	r := NewResettableRateCounter(interval)

	for i := 0; i < b.N; i++ {
		r.Incr(1)
		r.Rate()
	}
}

func BenchmarkResettableRateCounter_Parallel(b *testing.B) {
	interval := 1 * time.Millisecond
	r := NewResettableRateCounter(interval)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Incr(1)
			r.Rate()
		}
	})
}

func BenchmarkResettableRateCounter_With5MillionExisting(b *testing.B) {
	interval := 1 * time.Hour
	r := NewResettableRateCounter(interval)

	for i := 0; i < 5000000; i++ {
		r.Incr(1)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Incr(1)
		r.Rate()
	}
}
