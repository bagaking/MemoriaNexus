package utils

import (
	"context"
	"sync"
	"time"
)

// BatchProcessor defines a function to process a batch of tasks.
type BatchProcessor func(ctx context.Context, start, end int) error

// ParallelOptions defines options for parallel batch processing.
type ParallelOptions struct {
	MaxGoroutines int
	Timeout       time.Duration
	BatchSize     int
}

// Option defines a function to set options for parallel batch processing.
type Option func(*ParallelOptions)

// BatchWithMaxGoroutines sets the maximum number of concurrent goroutines.
func BatchWithMaxGoroutines(max int) Option {
	return func(opts *ParallelOptions) {
		opts.MaxGoroutines = max
	}
}

// BatchWithTimeout sets a timeout duration for batch processing.
func BatchWithTimeout(timeout time.Duration) Option {
	return func(opts *ParallelOptions) {
		opts.Timeout = timeout
	}
}

// BatchWithSize sets the size for each batch.
func BatchWithSize(size int) Option {
	return func(opts *ParallelOptions) {
		opts.BatchSize = size
	}
}

// ParallelBatchProcess processes tasks in parallel batches with specified options.
func ParallelBatchProcess(ctx context.Context, total int, processor BatchProcessor, opts ...Option) error {
	options := ParallelOptions{
		MaxGoroutines: 10,          // Default maximum number of concurrent goroutines
		Timeout:       time.Minute, // Default timeout duration
		BatchSize:     500,         // Default batch size
	}

	// Apply provided options to customize parallel processing behavior
	for _, opt := range opts {
		opt(&options)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, (total/options.BatchSize)+1)
	sem := make(chan struct{}, options.MaxGoroutines) // Semaphore to control concurrency
	ctx, cancel := context.WithTimeout(ctx, options.Timeout)
	defer cancel() // Ensure context cancellation when function exits

	// Process each batch of tasks in parallel
	for start := 0; start < total; start += options.BatchSize {
		end := start + options.BatchSize
		if end > total {
			end = total
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			select {
			case sem <- struct{}{}: // Acquire semaphore
				defer func() { <-sem }() // Release semaphore
				if err := processor(ctx, start, end); err != nil {
					errChan <- err // Send error to channel
				}
			case <-ctx.Done():
				errChan <- ctx.Err() // Send context error to channel
			}
		}(start, end)
	}

	wg.Wait()
	close(errChan)

	// Return the first error occurred
	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}
