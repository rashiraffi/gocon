# gocon: A Concurrency Utility Library for Go

**gocon** is a Go package that provides utilities for working with concurrency in Go, including worker pool functionality to simplify running tasks concurrently. It is designed to be flexible and easy to use, supporting context cancellation and error handling for effective concurrent processing.

## Features

- **Worker Pool without Error Handling**: Easily create a worker pool to process jobs concurrently.
- **Worker Pool with Error Handling**: Process jobs concurrently while also capturing and handling errors.
- **Context Support**: Use Go's context package to manage cancellation and timeouts for graceful shutdown of workers.
- **Highly Configurable**: Customize the number of workers, job processing function, and manage channels for results and errors.

## Installation

Add `gocon` to your project using:

```sh
go get github.com/yourusername/gocon
```

## Usage

The package provides two primary functions: `RunWorkerPool` and `RunWorkerPoolWithErrors`. Here are examples of how to use them.

### RunWorkerPool (Without Error Handling)

`RunWorkerPool` starts a worker pool to process jobs concurrently, using a job processing function. The worker pool is context-aware, which allows graceful shutdown when the context is canceled.

```go
import (
    "context"
    "fmt"
    "time"
    "github.com/yourusername/gocon"
)

func main() {
    jobs := make(chan int, 10)
    for i := 0; i < 10; i++ {
        jobs <- i
    }
    close(jobs)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    results := gocon.RunWorkerPool(ctx, jobs, 3, func(ctx context.Context, job int) string {
        return fmt.Sprintf("Processed job: %d", job)
    })

    for result := range results {
        fmt.Println(result)
    }
}
```

### RunWorkerPoolWithErrors (With Error Handling)

`RunWorkerPoolWithErrors` provides similar functionality to `RunWorkerPool`, but also captures and reports errors from the job processing function.

```go
import (
    "context"
    "fmt"
    "time"
    "github.com/yourusername/gocon"
)

func main() {
    jobs := make(chan int, 10)
    for i := 0; i < 10; i++ {
        jobs <- i
    }
    close(jobs)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    results, errors := gocon.RunWorkerPoolWithErrors(ctx, jobs, 3, func(ctx context.Context, job int) (string, error) {
        if job < 0 {
            return "", fmt.Errorf("negative job: %d", job)
        }
        return fmt.Sprintf("Processed job: %d", job), nil
    })

    // Handle errors concurrently
    go func() {
        for err := range errors {
            fmt.Printf("Error: %v\n", err)
        }
    }()

    for result := range results {
        fmt.Println(result)
    }
}
```

## API Reference

- `RunWorkerPool[In any, Out any](ctx context.Context, jobs <-chan In, workerCount int, workerFunc func(context.Context, In) Out) <-chan Out`

  - Starts a worker pool to process jobs concurrently without error handling.
  - **Parameters**:
    - `ctx`: Context for managing cancellation.
    - `jobs`: Channel of jobs to process.
    - `workerCount`: Number of workers to spawn.
    - `workerFunc`: Function to process each job.
  - **Returns**: Channel of results.

- `RunWorkerPoolWithErrors[In any, Out any](ctx context.Context, jobs <-chan In, workerCount int, workerFunc func(context.Context, In) (Out, error)) (<-chan Out, <-chan error)`
  - Starts a worker pool to process jobs concurrently with error handling.
  - **Parameters**:
    - `ctx`: Context for managing cancellation.
    - `jobs`: Channel of jobs to process.
    - `workerCount`: Number of workers to spawn.
    - `workerFunc`: Function to process each job, which may return an error.
  - **Returns**: Channel of results and channel of errors.

## Contributing

Contributions are welcome! Feel free to submit a pull request or file an issue if you have suggestions or improvements.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments

Thanks to the Go community for providing a great environment for learning and building concurrent systems.
