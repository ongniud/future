
# Future Package

`Future` is a Go package for asynchronous computation. It allows you to initiate an asynchronous task and retrieve its result later. `Future` supports lazy execution mode and task cancellation.

## Features

- **Asynchronous Computation**: Initiates asynchronous tasks and allows retrieving results when needed.
- **Lazy Execution Mode**: Supports lazy execution where the task is only run when `Await` is called for the first time.
- **Task Cancellation**: Allows cancelling tasks using the `Abort` method.
- **Thread Safety**: Ensures thread safety using `sync.Once` and `atomic.Value`.

## Installation

Install the package using `go get`:

```bash
go get github.com/ongniud/future
```

## Usage

### Creating a Future

To create a new Future, use the `NewFuture` function:

```go
f := future.NewFuture(ctx, promiseFunc, future.WithLazy())
```

You can pass in a context, a promise function, and options like `WithLazy` to control the execution behavior.

### Waiting for the Future to Complete

To wait for the Future to complete and retrieve the result, use the `Await` method:

```go
result, err := f.Await()
```

This method will block until the Future completes or the context is canceled.

### Cancelling the Future

To cancel the Future, use the `Abort` method:

```go
f.Abort()
```

This will cancel the Future and stop its execution.

### Handling Completion

You can also use the `Done` method to get a channel that is closed when the Future completes.

```go
<-f.Done()
```

## Example
### Basic

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ongniud/future"
)

func main() {
	ctx := context.Background()
	promise := func(ctx context.Context) (any, error) {
		time.Sleep(2 * time.Second)
		return "Hello, Future!", nil
	}

	// Create Future
	f := future.NewFuture(ctx, promise)

	// Retrieve result
	res, err := f.Await()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", res.(string)) // Output: Result: Hello, Future!
	}
}
```
### Lazy Execution Mode
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ongniud/future"
)

func main() {
	ctx := context.Background()
	promise := func(ctx context.Context) (any, error) {
		time.Sleep(2 * time.Second)
		return "Lazy Result", nil
	}

	// Create Future in lazy execution mode
	f := future.NewFuture(ctx, promise, future.WithLazy())

	// Retrieve result (task will execute at this point)
	res, err := f.Await()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", res.(string)) // Output: Result: Lazy Result
	}
}
```
### Task Cancellation

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ongniud/future"
)

func main() {
	ctx := context.Background()
	promise := func(ctx context.Context) (any, error) {
		time.Sleep(5 * time.Second)
		return "This will not be reached", nil
	}

	// Create Future
	f := future.NewFuture(ctx, promise)

	// Cancel the task after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		f.Abort()
	}()

	// Retrieve result
	res, err := f.Await()
	if err != nil {
		fmt.Println("Error:", err) // Output: Error: context canceled
	} else {
		fmt.Println("Result:", res)
	}
}

```

## API Documentation
### NewFuture

```go 
func NewFuture(ctx context.Context, promise func(context.Context) (any, error), opts ...Option) *Future
```

Creates a new Future instance.
- ctx: Context used for task cancellation and timeouts.
- promise: The asynchronous task function.
- opts: Optional parameters like WithLazy.

### WithLazy

```go
func WithLazy() Option
```

Enables lazy execution mode. The task will execute only when Await is called for the first time.

### Await
```go
func (f *Future) Await() (any, error)
```

Waits for the task to complete and returns the result. If the task has not been started (lazy mode), it will trigger the task.

### Done

```go
func (f *Future) Done() <-chan struct{}
```

Returns a channel that will be closed when the task completes.

### Abort

```go
func (f *Future) Abort()
```
Cancels the task.

## Design Details
- Thread Safety: Uses sync.Once to ensure the task is executed only once, and atomic.Value to store the result and error.
- Lazy Execution Mode: Enabled through the WithLazy option, the task is started when Await is called for the first time.
- Task Cancellation: Task cancellation is implemented through context.Context.

## Contribution
Feel free to submit issues and pull requests! Please ensure code style consistency and passing tests.

## License
This project is licensed under the MIT License.
