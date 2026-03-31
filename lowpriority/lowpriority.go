// Package lowpriority provides a way to run code with low priority.
package lowpriority

import (
	"context"
	"runtime"
	"syscall"

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
		runtime.LockOSThread()
		tid := syscall.Gettid()
		err := syscall.Setpriority(syscall.PRIO_PROCESS, tid, 19)
		if err != nil {
			panic(err)
		}
		f(ctx)
	}).Wait()
}
