// Package lowpriority provides a way to run code with low priority.
package lowpriority

import (
	"context"
	"runtime"
	"sync"

	"github.com/pierrre/go-libs/funcutil"
	"github.com/pierrre/go-libs/goroutine"
)

// Run runs the given function with a low priority thread.
//
// The function is executed in a new goroutine which is locked to an OS thread.
// The thread priority is set to the lowest possible value (19 on Linux).
// The thread is terminated after the function returns, because it's not possible to restore the original priority without special privileges.
// The priority is not inherited by goroutines started by the function.
// [Run] waits for the function to finish, and returns any panic that may have occurred in the function.
func Run(ctx context.Context, f func(ctx context.Context)) {
	goroutine.Start(ctx, func(ctx context.Context) {
		setLowPriority()
		f(ctx)
	}).Wait()
}

// Pool creates a pool of workers that runs functions with low priority.
func Pool(ctx context.Context, workers int) (run func(ctx context.Context, f func(ctx context.Context)) error, stop func()) {
	type result struct {
		goexit   bool
		panicErr error
	}
	type task struct {
		ctx   context.Context //nolint:containedctx // The context is used to run functions.
		f     func(ctx context.Context)
		resCh chan<- result
	}
	taskCh := make(chan task)
	waiter := goroutine.StartN(ctx, workers, func(ctx context.Context, i int) {
		setLowPriority()
		for t := range taskCh {
			func() {
				var res result
				defer func() {
					t.resCh <- res
				}()
				funcutil.Call(
					func() {
						t.f(t.ctx)
					},
					func(goexit bool, panicErr error) {
						res.goexit = goexit
						res.panicErr = panicErr
					},
				)
			}()
		}
	})
	run = func(ctx context.Context, f func(ctx context.Context)) error {
		resCh := make(chan result, 1)
		t := task{
			ctx:   ctx,
			f:     f,
			resCh: resCh,
		}
		select {
		case taskCh <- t:
		case <-ctx.Done():
			return context.Cause(ctx)
		}
		var res result
		select {
		case res = <-resCh:
		case <-ctx.Done():
			return context.Cause(ctx)
		}
		if res.goexit {
			runtime.Goexit()
		}
		if res.panicErr != nil {
			panic(res.panicErr)
		}
		return nil
	}
	stop = sync.OnceFunc(func() {
		close(taskCh)
		waiter.Wait()
	})
	return run, stop
}

func setLowPriority() {
	runtime.LockOSThread()
	err := setThreadLowPriority()
	if err != nil {
		panic(err)
	}
}
