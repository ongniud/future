package future

import (
	"context"
	"sync"
	"sync/atomic"
)

// Option defines functional options for Future.
type Option func(*Future)

// WithLazy enables lazy execution of the Future.
func WithLazy() Option {
	return func(f *Future) {
		f.lazy = true
	}
}

// Future represents an asynchronous computation.
type Future struct {
	ctx     context.Context
	cancel  context.CancelFunc
	promise func(ctx context.Context) (any, error)

	// options
	lazy bool

	// result
	done chan struct{}
	res  atomic.Value
	err  atomic.Value
	once sync.Once
}

// NewFuture creates a new Future.
func NewFuture(ctx context.Context, promise func(context.Context) (any, error), opts ...Option) *Future {
	newCtx, cancel := context.WithCancel(ctx)
	f := &Future{
		ctx:     newCtx,
		cancel:  cancel,
		promise: promise,
		done:    make(chan struct{}),
	}

	// Apply options
	for _, opt := range opts {
		opt(f)
	}

	// Start execution if not lazy
	if !f.lazy {
		go f.run()
	}

	return f
}

// run executes the promise and stores the result.
func (f *Future) run() {
	f.once.Do(func() {
		res, err := f.promise(f.ctx)
		if res != nil {
			f.res.Store(res)
		}
		if err != nil {
			f.err.Store(err)
		}
		close(f.done) // Ensure this is called only once
	})
}

// Await waits for the Future to complete and returns the result.
func (f *Future) Await() (any, error) {
	// Trigger execution for lazy mode
	if f.lazy {
		go f.run()
	}

	// Wait for completion or context cancellation
	select {
	case <-f.done:
		res := f.res.Load()
		err, _ := f.err.Load().(error)
		return res, err
	case <-f.ctx.Done():
		if err := f.ctx.Err(); err != nil {
			return nil, err
		}
		return nil, context.Canceled
	}
}

// Done returns a channel that is closed when the Future completes.
func (f *Future) Done() <-chan struct{} {
	return f.done
}

// Abort cancels the Future.
func (f *Future) Abort() {
	f.cancel()
}
