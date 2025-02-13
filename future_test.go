package future

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestFuture_Success(t *testing.T) {
	f := NewFuture(context.Background(), func(ctx context.Context) (any, error) {
		return 42, nil
	})

	res, err := f.Await()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}
}

func TestFuture_LazyExecution(t *testing.T) {
	executed := false
	f := NewFuture(context.Background(), func(ctx context.Context) (any, error) {
		executed = true
		return 100, nil
	}, WithLazy())

	time.Sleep(100 * time.Millisecond) // 确保 lazy 模式不会立即执行
	if executed {
		t.Errorf("Future executed before Await() was called")
	}

	res, err := f.Await()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res != 100 {
		t.Errorf("Expected 100, got %v", res)
	}
}

func TestFuture_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	f := NewFuture(ctx, func(ctx context.Context) (any, error) {
		time.Sleep(200 * time.Millisecond) // 超时
		return "hello", nil
	})

	res, err := f.Await()
	if err == nil {
		t.Errorf("Expected timeout error, got nil")
	}
	if res != nil {
		t.Errorf("Expected nil result, got %v", res)
	}
}

func TestFuture_AbortBeforeRun(t *testing.T) {
	ctx := context.Background()
	f := NewFuture(ctx, func(ctx context.Context) (any, error) {
		time.Sleep(500 * time.Millisecond)
		return "done", nil
	})

	f.Abort()

	res, err := f.Await()
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
	if res != nil {
		t.Errorf("Expected nil result, got %v", res)
	}
}

func TestFuture_AbortAfterRun(t *testing.T) {
	ctx := context.Background()
	f := NewFuture(ctx, func(ctx context.Context) (any, error) {
		return "finished", nil
	})

	time.Sleep(50 * time.Millisecond) // 确保任务执行
	f.Abort()

	res, err := f.Await()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res != "finished" {
		t.Errorf("Expected 'finished', got %v", res)
	}
}

func TestFuture_ConcurrentAwait(t *testing.T) {
	f := NewFuture(context.Background(), func(ctx context.Context) (any, error) {
		time.Sleep(100 * time.Millisecond)
		return "hello", nil
	})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		res, err := f.Await()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res != "hello" {
			t.Errorf("Expected 'hello', got %v", res)
		}
	}()

	go func() {
		defer wg.Done()
		res, err := f.Await()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res != "hello" {
			t.Errorf("Expected 'hello', got %v", res)
		}
	}()

	wg.Wait()
}
