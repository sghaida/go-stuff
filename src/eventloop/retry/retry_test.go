package retry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry_Run(t *testing.T) {
	// fake error
	var testErr = errors.New("test error")

	t.Run("r.Run returns nothing with error", func(t *testing.T) {
		tries := 0
		start := time.Now()
		var last time.Time

		retry := NewRetry(5, 50*time.Millisecond, 50*time.Millisecond)

		res, err := retry.Run(func() (interface{}, error) {
			tries++
			last = time.Now()
			return nil, testErr
		})

		assert.Equal(t, tries, 5, fmt.Sprintf("expected 5 tries, got %d", tries))
		assert.Equal(t, err, testErr, fmt.Sprintf("err should equal testErr, got: %v", err))
		assert.Equal(t, res, nil)

		// CI environments and VMs are unreliable when it comes to firing
		// time.Timer, so this assertion is skipped on CI, because it was
		// causing random failures.
		if _, present := os.LookupEnv("CI"); present {
			max := 5 * (50 + 50) * time.Millisecond
			assert.Less(
				t, last.Sub(start).Milliseconds(), max.Milliseconds(),
				fmt.Sprintf("should have taken less than %v, took %v", max.Milliseconds(), last.Sub(start).Milliseconds()),
			)
		}
	})

	t.Run("r.Run returns something after three retries", func(t *testing.T) {
		tries := 0
		retry := NewRetry(5, 50*time.Millisecond, 50*time.Millisecond)

		res, err := retry.Run(func() (interface{}, error) {
			if tries == 3 {
				return tries, nil
			}
			tries++
			return nil, testErr
		})

		assert.Equal(t, tries, 3, fmt.Sprintf("expected 3 tries, got %d", tries))
		assert.Equal(t, err, nil, fmt.Sprintf("err should equal testErr, got: %v", err))
		assert.Equal(t, res, 3)
	})
}

func TestRetry_RunWithContext(t *testing.T) {
	// fake error
	var testErr = errors.New("test error")

	t.Run("r.RunWithContext proper timeout", func(t *testing.T) {
		tries := 0
		start := time.Now()
		var last time.Time

		retry := NewRetry(5, 50*time.Millisecond, 50*time.Millisecond)

		ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
		defer cancel()

		res, err := retry.RunWithContext(ctx, func(ctx context.Context) (interface{}, error) {
			tries++
			last = time.Now()
			return nil, testErr
		})

		assert.Equal(t, tries, 5, fmt.Sprintf("expected 5 tries, got %d", tries))
		assert.Equal(t, err, testErr, fmt.Sprintf("err should equal errTest, got: %v", err))
		assert.Equal(t, res, nil)

		// CI environments and VMs are unreliable when it comes to firing
		// time.Timer, so this assertion is skipped on CI, because it was
		// causing random failures.
		if _, present := os.LookupEnv("CI"); present {
			max := 5 * (50 + 50) * time.Millisecond
			assert.Less(
				t, last.Sub(start).Milliseconds(), max.Milliseconds(),
				fmt.Sprintf("should have taken less than %v, took %v", max.Milliseconds(), last.Sub(start).Milliseconds()),
			)
		}
	})

	t.Run("r.RunWithContext small timeout", func(t *testing.T) {
		tries := 0
		start := time.Now()
		var last time.Time

		retry := NewRetry(5, 50*time.Millisecond, 50*time.Millisecond)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		res, err := retry.RunWithContext(ctx, func(ctx context.Context) (interface{}, error) {
			tries++
			last = time.Now()
			return nil, testErr
		})

		assert.Less(t, tries, 5, fmt.Sprintf("expected less than 5 tries, got %d", tries))
		assert.Equal(t, err, testErr, fmt.Sprintf("err should equal errTest, got: %v", err))
		assert.Equal(t, res, nil)

		// CI environments and VMs are unreliable when it comes to firing
		// time.Timer, so this assertion is skipped on CI, because it was
		// causing random failures.
		if _, present := os.LookupEnv("CI"); present {
			max := 5 * (50 + 50) * time.Millisecond
			assert.Less(
				t, last.Sub(start).Milliseconds(), max.Milliseconds(),
				fmt.Sprintf("should have taken less than %v, took %v", max.Milliseconds(), last.Sub(start).Milliseconds()),
			)
		}
	})

	t.Run("exit due to context cancellation", func(t *testing.T) {
		var wg sync.WaitGroup
		var err error

		tries := 0
		retry := NewRetry(5, 50*time.Millisecond, 50*time.Millisecond)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

		wg.Add(1)
		go func() {
			_, err = retry.RunWithContext(ctx, func(ctx context.Context) (interface{}, error) {
				tries++
				return nil, testErr
			})
			wg.Done()
		}()

		time.Sleep(200 * time.Millisecond)
		cancel()
		wg.Wait()

		assert.GreaterOrEqual(t, tries, 1, fmt.Sprintf("expected at least 1, got %d", tries))
		assert.Equal(t, err, testErr, fmt.Sprintf("err should equal Test error, got: %v", err))
	})

	t.Run("less than zero durations", func(t *testing.T) {
		var tries int
		retry := NewRetry(-1, -1, -1)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err := retry.RunWithContext(ctx, func(ctx context.Context) (interface{}, error) {
			tries++
			return nil, testErr
		})

		assert.GreaterOrEqual(t, tries, 1, fmt.Sprintf("expected at less than 3, got %d", tries))
		assert.Equal(t, err, testErr, fmt.Sprintf("err should equal Test error, got: %v", err))
	})
}
