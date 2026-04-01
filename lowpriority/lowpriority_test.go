package lowpriority

import (
	"context"
	"crypto/sha512"
	"runtime"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/go-libs/goroutine"
)

var data = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

func TestRun(t *testing.T) {
	ctx := t.Context()
	Run(ctx, func(ctx context.Context) {
		task()
	})
}

func BenchmarkRunNormalPriority(b *testing.B) {
	defer startConsumeCPU()()
	for b.Loop() {
		task()
	}
}

func BenchmarkRunLowPriority(b *testing.B) {
	ctx := b.Context()
	defer startConsumeCPU()()
	Run(ctx, func(ctx context.Context) {
		for b.Loop() {
			task()
		}
	})
}

func TestPool(t *testing.T) {
	ctx := t.Context()
	run, stop := Pool(ctx, 1)
	defer stop()
	err := run(ctx, func(ctx context.Context) {
		task()
	})
	assert.NoError(t, err)
}

func BenchmarkPool(b *testing.B) {
	ctx := b.Context()
	run, stop := Pool(ctx, 1)
	defer stop()
	for b.Loop() {
		err := run(ctx, func(ctx context.Context) {
			task()
		})
		assert.NoError(b, err)
	}
}

func startConsumeCPU() (stop func()) {
	return goroutine.StartWithCancel(context.Background(), consumeCPU).Wait
}

func consumeCPU(ctx context.Context) {
	goroutine.RunN(ctx, runtime.GOMAXPROCS(0), func(ctx context.Context, i int) {
		for ctx.Err() == nil {
			task()
		}
	})
}

func task() {
	sha512.Sum512(data)
}
