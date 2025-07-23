package utils

import (
	"errors"
	"log"
	"time"
)

// RetryFunc is the function signature for retryable operations.
type RetryFunc[T any] func() (T, error)

// Retry retries a function up to maxAttempts times with delay between attempts.
// If the function succeeds (err == nil), it returns the result.
// If all attempts fail, the last error is returned.
func Retry[T any](maxAttempts int, delay time.Duration, fn RetryFunc[T]) (T, error) {
	var lastErr error
	var result T

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, lastErr = fn()
		if lastErr == nil {
			return result, nil
		}

		log.Printf("retry attempt %d failed: %v", attempt, lastErr)

		if attempt < maxAttempts {
			time.Sleep(delay)
		}
	}

	return result, errors.New("all retry attempts failed: " + lastErr.Error())
}