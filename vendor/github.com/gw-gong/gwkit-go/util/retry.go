package util

import (
	"context"
	"fmt"
	"time"
)

const (
	DefaultMaxTries      = 3
	DefaultRetryInterval = 10 * time.Millisecond
	MaxBackoff           = 10 * time.Second
)

// WithRetry is a function that retries a function until it returns false or the context is canceled.
// maxTries is the maximum number of tries, interval is the interval between tries.
// allTries is the number of tries, allErrs is the "returned error" of all tries, lastErr is the error of the last try.
func WithRetry(ctx context.Context, maxTries int, interval time.Duration, operationFunc func() (shouldRetry bool, err error)) (allTries int, allErrs []error, lastErr error) {
	if ctx == nil {
		return 0, []error{fmt.Errorf("context cannot be nil")}, fmt.Errorf("context cannot be nil")
	}
	if operationFunc == nil {
		return 0, []error{fmt.Errorf("operationFunc cannot be nil")}, fmt.Errorf("operationFunc cannot be nil")
	}
	if maxTries <= 0 {
		maxTries = DefaultMaxTries
	}
	if interval <= 0 {
		interval = DefaultRetryInterval
	}

	allErrs = make([]error, 0, maxTries)
	for i := 1; i <= maxTries; i++ {
		shouldRetry, err := operationFunc()
		lastErr = err
		if err != nil {
			allErrs = append(allErrs, fmt.Errorf("try %d/%d failed: %w", i, maxTries, lastErr))
		}
		if !shouldRetry {
			return i, allErrs, lastErr
		}

		// The last loop does not need to wait, just return
		if i == maxTries {
			break
		}

		select {
		case <-ctx.Done():
			lastErr = fmt.Errorf("canceled after %d/%d try: %w", i, maxTries, ctx.Err())
			allErrs = append(allErrs, lastErr)
			return i, allErrs, lastErr
		default:
			// Exponential backoff: 1x, 2x, 4x, 8x, ...
			backoff := time.Duration(1<<uint(i-1)) * interval
			if backoff > MaxBackoff {
				backoff = MaxBackoff
			}
			time.Sleep(backoff)
		}
	}
	return maxTries, allErrs, lastErr
}
