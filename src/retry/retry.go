package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

const (
	DefaultMaxTries     = 5
	DefaultInitialDelay = time.Millisecond * 200
	DefaultMaxDelay     = time.Millisecond * 1000
)

// Retry contains the max retires and min|max delay
type Retry struct {
	maxTries     int
	initialDelay time.Duration
	maxDelay     time.Duration
}

// NewRetry initialize new retry
func NewRetry(maxTries int, initialDelay, maxDelay time.Duration) *Retry {
	if maxTries <= 0 {
		maxTries = DefaultMaxTries
	}
	if initialDelay <= 0 {
		initialDelay = DefaultInitialDelay
	}
	if maxDelay <= 0 {
		maxDelay = DefaultMaxDelay
	}
	return &Retry{maxTries, initialDelay, maxDelay}
}

// Run runs the function that needs to be retried until it returns nil or termination error or the context is done
// Run should return results as interface and error
func (r *Retry) Run(funcToRetry func() (interface{}, error)) (interface{}, error) {
	return r.RunWithContext(context.Background(), func(_ context.Context) (result interface{}, err error) {
		result, err = funcToRetry()
		return
	})
}

// RunWithContext runs the function that needs to be retried until it returns nil or termination error or the context is done
// the retry is going to work with exponential backoff until the context is deadlined
// execution return some result interface or error
func (r *Retry) RunWithContext(
	ctx context.Context, funcToRetry func(context.Context) (interface{}, error),
) (interface{}, error) {
	// create a random source which is used in setting the jitter

	attempts := 0
	for {
		// run the function
		result, err := funcToRetry(ctx)
		// no error, then we are done!
		if err == nil {
			return result, nil
		}
		// max retries is reached return the error
		attempts++
		if attempts == r.maxTries {
			return nil, err
		}
		// wait for the next duration or context canceled|time out, whichever comes first
		t := time.NewTimer(getNextBackoff(attempts, r.initialDelay, r.maxDelay))
		select {
		case <-t.C:
			// nothing to be done as the timer is killed
		case <-ctx.Done():
			// context cancelled, kill the timer if it is not killed, and return the last error
			if !t.Stop() {
				<-t.C
			}
			return nil, err
		}
	}
}

// getNextBackoff based on https://en.wikipedia.org/wiki/Exponential_backoff
func getNextBackoff(attempts int, initialDelay, maxDelay time.Duration) time.Duration {

	mx := float64(maxDelay)
	mn := float64(initialDelay)

	dur := mn * math.Pow(2, float64(attempts))
	if dur > mx {
		dur = mx
	}

	j := time.Duration(rand.Float64()*(dur-mn) + mn)
	return j
}
