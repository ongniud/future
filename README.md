# Future

The Future package provides a simple and efficient way to handle asynchronous computations in Go. It allows you to execute tasks concurrently, retrieve their results, and cancel them if needed. The package supports lazy execution, panic recovery, and context-based cancellation.

## Features

- **Asynchronous Execution:** Execute tasks concurrently and retrieve their results when needed.
- **Lazy Execution:** Delay task execution until the result is explicitly requested.
- **Task Cancellation:** Cancel tasks manually using the `Abort()` method.
- **Panic Recovery:** Automatically recover from panics in tasks and return them as errors.
- **Thread Safety:** Safe for concurrent use with proper synchronization.

## Installation

To use the Future package, simply import it in your Go project:

```go
import "github.com/ongniud/future"
```

## Usage

### Creating a Future

Use the `NewFuture` function to create a new Future:

```go
task := func(ctx context.Context) (any, error) {
    // Simulate a long-running task
    time.Sleep(100 * time.Millisecond)
    return "task completed", nil
}

f := A.NewFuture(context.Background(), task)
```

### Retrieving the Result

Use the `Result()` method to wait for the task to complete and retrieve its result:

```go
res, err := f.Result()
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Result:", res)
}
```

### Lazy Execution

To enable lazy execution, use the `WithLazy` option:

```go
f := A.NewFuture(context.Background(), task, A.WithLazy())
```

In lazy mode, the task will only start when `Result()` is called.

### Canceling a Task

Use the `Abort()` method to cancel a task:

```go
f.Abort()
```

If the task is canceled, `Result()` will return a `context.Canceled` error.

### Checking Task Status

Use the `Ready()` method to check if the task has completed:

```go
if f.Ready() {
    fmt.Println("Task is ready")
} else {
    fmt.Println("Task is still running")
}
```

### Waiting for Completion

Use the `Done()` method to get a channel that is closed when the task completes:

```go
<-f.Done()
fmt.Println("Task is done")
```

### Example

Here is a complete example demonstrating the usage of the Future package:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/ongniud/future"
)

func main() {
    task := func(ctx context.Context) (any, error) {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(200 * time.Millisecond):
            return "task completed", nil
        }
    }

    f := A.NewFuture(context.Background(), task)

    // Simulate cancellation after 100ms
    go func() {
        time.Sleep(100 * time.Millisecond)
        f.Abort()
    }()

    // Wait for the result
    res, err := f.Result()
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", res)
    }
}
```
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
