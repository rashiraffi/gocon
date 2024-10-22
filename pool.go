package gocon

import (
	"context"
	"sync"
)

// RunWorkerPool starts a worker pool to process jobs concurrently, without error handling.
// This function spawns workerCount number of goroutines to process the jobs provided via the jobs channel.
// Each worker uses the workerFunc to process a job, and the results are sent to the result channel.
func RunWorkerPool[In any, Out any](ctx context.Context, jobs <-chan In, workerCount int, workerFunc func(context.Context, In) Out) <-chan Out {
	result := make(chan Out) // Channel to collect results from the workers.
	var wg sync.WaitGroup

	// Add workerCount to the WaitGroup to ensure we wait for all workers to finish.
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done() // Decrement the WaitGroup counter when the worker is done.
			for job := range jobs {
				select {
				case <-ctx.Done():
					// If the context is canceled, exit the worker.
					return
				default:
					// Process the job using the provided worker function.
					output := workerFunc(ctx, job)
					// Send the result to the result channel.
					result <- output
				}
			}
		}()
	}

	// Close the result channel once all workers are done.
	go func() {
		wg.Wait()     // Wait for all workers to complete.
		close(result) // Close the result channel to signal that no more results will be sent.
	}()

	return result
}

// RunWorkerPoolWithErrors starts a worker pool to process jobs concurrently, with error handling.
// This function works similarly to RunWorkerPool, but also captures errors from the workerFunc.
func RunWorkerPoolWithErrors[In any, Out any](ctx context.Context, jobs <-chan In, workerCount int, workerFunc func(context.Context, In) (Out, error)) (<-chan Out, <-chan error) {
	result := make(chan Out)   // Channel to collect results from the workers.
	errors := make(chan error) // Channel to collect errors from the workers.
	var wg sync.WaitGroup

	// Add workerCount to the WaitGroup to ensure we wait for all workers to finish.
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done() // Decrement the WaitGroup counter when the worker is done.
			for job := range jobs {
				select {
				case <-ctx.Done():
					// If the context is canceled, exit the worker.
					return
				default:
					// Process the job using the provided worker function, capturing any error.
					output, err := workerFunc(ctx, job)
					if err != nil {
						// Send the error to the errors channel if one occurs.
						errors <- err
						continue
					}
					// Send the result to the result channel if no error occurred.
					result <- output
				}
			}
		}()
	}

	// Close the result and error channels once all workers are done.
	go func() {
		wg.Wait()     // Wait for all workers to complete.
		close(result) // Close the result channel to signal that no more results will be sent.
		close(errors) // Close the errors channel to signal that no more errors will be sent.
	}()

	return result, errors
}
