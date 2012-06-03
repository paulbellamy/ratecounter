# ratecounter

A Concurrent RateCounter implementation in Golang

## Usage

```
import "github.com/paulbellamy/ratecounter"
```

Package ratecounter provides a concurrent rate-counter, for tracking
counts in an interval

## Documentation

```
type Counter struct {
    // contains filtered or unexported fields
}
    A Counter is a thread-safe counter implementation

func NewCounter() *Counter
    NewCounter is used to construct a new Counter object

func (c *Counter) Incr(val int64)
    Increment the counter by some value

func (c *Counter) Value() int64
    Return the counter's current value

type RateCounter struct {
    // contains filtered or unexported fields
}
    A RateCounter is a thread-safe counter which returns the number of times
    'Mark' has been called in the last interval

func NewRateCounter(intrvl time.Duration) *RateCounter
    Constructs a new RateCounter, for the interval provided

func (r *RateCounter) Mark()
    Add 1 event into the RateCounter

func (r *RateCounter) Rate() int64
    Return the current number of events in the last interval
```
