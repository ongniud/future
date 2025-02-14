package A

import (
	"context"
	"fmt"
	"sync"
)

// Option defines functional options for Future.
type Option func(*Future)

// WithLazy enables lazy execution of the Future.
func WithLazy() Option {
	return func(f *Future) {
		f.lazy = true
	}
}

type Future struct {
	task func(context.Context) (any, error)
	lazy bool

	item interface{}
	err  error

	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
	once   sync.Once
	closed sync.Once
	done   chan struct{}
}

// Result waits for the result to be ready and returns it.
func (f *Future) Result() (interface{}, error) {
	f.once.Do(f.start)
	<-f.done

	f.mu.Lock()
	defer f.mu.Unlock()
	return f.item, f.err
}

// Ready returns true if the result is available.
func (f *Future) Ready() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

// Done returns a channel that is closed when the result is ready.
func (f *Future) Done() <-chan struct{} {
	return f.done
}

// NewFuture creates a new Future.
func NewFuture(ctx context.Context, task func(context.Context) (any, error), opts ...Option) *Future {
	newCtx, cancel := context.WithCancel(ctx)
	f := &Future{
		ctx:    newCtx,
		cancel: cancel,
		task:   task,
		done:   make(chan struct{}),
	}
	for _, opt := range opts {
		opt(f)
	}
	if !f.lazy {
		f.once.Do(f.start)
	}
	return f
}

// start executes the task and stores the result.
func (f *Future) start() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.mu.Lock()
				f.err = fmt.Errorf("panic occurred: %v", r)
				f.mu.Unlock()
				f.markDone()
			}
		}()
		res, err := f.task(f.ctx)
		f.mu.Lock()
		f.item, f.err = res, err
		f.mu.Unlock()
		f.markDone()
	}()
}

// Abort cancels the task execution.
func (f *Future) Abort() {
	f.mu.Lock()
	defer f.mu.Unlock()
	select {
	case <-f.done:
		return
	default:
		if f.cancel != nil {
			f.cancel()
		}
		f.item, f.err = nil, context.Canceled
		f.markDone()
	}
}

// markDone marks the future as done and closes the done channel.
func (f *Future) markDone() {
	f.closed.Do(func() {
		close(f.done)
	})
}
