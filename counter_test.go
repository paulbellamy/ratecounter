package ratecounter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	var c Counter

	assert.Equal(t, int64(0), c.Value())
	c.Incr(1)
	assert.Equal(t, int64(1), c.Value())
	c.Incr(9)
	assert.Equal(t, int64(10), c.Value())
	c.Reset()
	assert.Equal(t, int64(0), c.Value())

	// Concurrent usage
	wg := &sync.WaitGroup{}
	wg.Add(3)
	for i := 1; i <= 3; i++ {
		go func(val int64) {
			c.Incr(val)
			wg.Done()
		}(int64(i))
	}
	wg.Wait()
	assert.Equal(t, int64(6), c.Value())
}

func BenchmarkCounter(b *testing.B) {
	var c Counter

	for i := 0; i < b.N; i++ {
		c.Incr(1)
	}
}
