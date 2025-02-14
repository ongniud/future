package A

import (
	"context"
	"testing"
	"time"
)

func TestFuture_Result(t *testing.T) {
	task := func(ctx context.Context) (any, error) {
		return "success", nil
	}
	future := NewFuture(context.Background(), task)

	// Calling Result should wait for the task to complete and return the result
	result, err := future.Result()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "success" {
		t.Fatalf("expected result 'success', got %v", result)
	}
}

func TestFuture_Abort(t *testing.T) {
	task := func(ctx context.Context) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second): // Simulate a long-running task
			return "success", nil
		}
	}

	future := NewFuture(context.Background(), task)

	// Abort the task before it finishes
	future.Abort()

	// Call Result to ensure the task was canceled
	result, err := future.Result()
	if err == nil {
		t.Fatalf("expected error, got result %v", result)
	}
	if err.Error() != "context canceled" {
		t.Fatalf("expected 'context canceled' error, got %v", err)
	}
}

func TestFuture_PanicRecovery(t *testing.T) {
	task := func(ctx context.Context) (any, error) {
		panic("unexpected error")
	}

	future := NewFuture(context.Background(), task)

	// Calling Result should capture the panic and return it as an error
	_, err := future.Result()
	if err == nil {
		t.Fatalf("expected panic error, got nil")
	}
	if err.Error() != "panic occurred: unexpected error" {
		t.Fatalf("expected 'panic occurred: unexpected error', got %v", err)
	}
}

func TestFuture_LazyExecution(t *testing.T) {
	task := func(ctx context.Context) (any, error) {
		// Simulate work
		time.Sleep(500 * time.Millisecond)
		return "lazy success", nil
	}

	// Create future with lazy execution
	future := NewFuture(context.Background(), task, WithLazy())

	// `Result` should not block until it is called (lazy mode)
	start := time.Now()
	_, err := future.Result()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// The total duration should be more than 500ms to ensure that the task was lazy-loaded
	if duration < 500*time.Millisecond {
		t.Fatalf("expected lazy execution, took too little time %v", duration)
	}
}

func TestFuture_Ready(t *testing.T) {
	task := func(ctx context.Context) (any, error) {
		time.Sleep(100 * time.Millisecond)
		return "ready success", nil
	}
	future := NewFuture(context.Background(), task)

	// Initially, the task is not ready
	if future.Ready() {
		t.Fatalf("expected task to not be ready")
	}

	// Calling Result will make the task complete
	_, err := future.Result()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// After calling Result, the task should be ready
	if !future.Ready() {
		t.Fatalf("expected task to be ready")
	}
}

func TestFuture_DoneChannel(t *testing.T) {
	task := func(ctx context.Context) (any, error) {
		time.Sleep(100 * time.Millisecond)
		return "done success", nil
	}

	future := NewFuture(context.Background(), task)

	doneChan := future.Done()

	// Initially, done channel should not be closed
	select {
	case <-doneChan:
		t.Fatal("expected done channel to not be closed initially")
	default:
	}

	// Calling Result will complete the task and close the done channel
	_, err := future.Result()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Now the done channel should be closed
	select {
	case <-doneChan:
		// Done channel closed as expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected done channel to be closed")
	}
}
