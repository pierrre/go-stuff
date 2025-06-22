package funcmock

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/go-libs/goroutine"
)

func Example() {
	fs := []func(s string) int{
		func(s string) int {
			fmt.Println("function 1:", s)
			return 1
		},
		func(s string) int {
			fmt.Println("function 2:", s)
			return 2
		},
		func(s string) int {
			fmt.Println("function 3:", s)
			return 1
		},
	}
	onTooManyCalls := func(s string) int {
		fmt.Println("too many calls:", s)
		return 0
	}
	f := MockList(fs, onTooManyCalls)
	f, getCount := MockCount(f)
	f = MockSerial(f)
	fmt.Println(f("a"))
	fmt.Println(f("b"))
	fmt.Println(f("c"))
	fmt.Println(f("d"))
	fmt.Println(f("e"))
	fmt.Println("count:", getCount())
	// Output:
	// function 1: a
	// 1
	// function 2: b
	// 2
	// function 3: c
	// 1
	// too many calls: d
	// 0
	// too many calls: e
	// 0
	// count: 5
}

func TestMockList(t *testing.T) {
	onTooManyCallsCalled := false
	f := MockList([]func(s string) int{
		func(s string) int {
			assert.Equal(t, s, "aaa")
			return 1
		},
		func(s string) int {
			assert.Equal(t, s, "bbb")
			return 2
		},
	}, func(s string) int {
		onTooManyCallsCalled = true
		assert.Equal(t, s, "ccc")
		return 666
	})
	assert.Equal(t, f("aaa"), 1)
	assert.Equal(t, f("bbb"), 2)
	assert.Equal(t, f("ccc"), 666)
	assert.True(t, onTooManyCallsCalled)
}

func TestMockListPanicTooManyCalls(t *testing.T) {
	f := MockList([]func(){
		func() {},
	}, nil)
	f()
	assert.Panics(t, func() {
		f()
	})
}

func BenchmarkMockList(b *testing.B) {
	f := MockList(slices.Repeat([]func(){func() {}}, 100), func() {})
	for b.Loop() {
		f()
	}
}

func TestMockCount(t *testing.T) {
	f, getCount := MockCount(func() {})
	count := 10
	counts := make([]int64, 0, count)
	for range count {
		counts = append(counts, getCount())
		f()
	}
	assert.Equal(t, getCount(), int64(count))
	assert.SliceLen(t, counts, count)
	for i, n := range counts {
		assert.Equal(t, n, int64(i))
	}
}

func BenchmarkMockCount(b *testing.B) {
	f, _ := MockCount(func() {})
	for b.Loop() {
		f()
	}
}

func TestMockSerial(t *testing.T) {
	ctx := t.Context()
	count := 0
	f := MockSerial(func() {
		count++
	})
	workers := 8
	callsPerWorker := 100
	goroutine.RunN(ctx, workers, func(ctx context.Context) {
		for range callsPerWorker {
			f()
		}
	})
	assert.Equal(t, count, workers*callsPerWorker)
}

func BenchmarkMockSerial(b *testing.B) {
	f := MockSerial(func() {})
	for b.Loop() {
		f()
	}
}

func TestVariadic(t *testing.T) {
	f := mockTest(func(i int, vs ...string) {
		assert.Equal(t, i, 123)
		assert.SliceEqual(t, vs, []string{"a", "b", "c"})
	})
	f(123, "a", "b", "c")
}

func TestPanicKindNotFunc(t *testing.T) {
	assert.Panics(t, func() {
		mockTest("invalid")
	})
}

func BenchmarkMock(b *testing.B) {
	f := mockTest(func() {})
	for b.Loop() {
		f()
	}
}

func mockTest[F Func](f F) F {
	v := reflect.ValueOf(f)
	return MakeFunc[F](func(args []reflect.Value) []reflect.Value {
		return call(v, args)
	})
}

func TestMakeFuncResults(t *testing.T) {
	type ft func() (int, string)
	f := MakeFunc[ft](func(args []reflect.Value) []reflect.Value {
		return MakeFuncResults[ft]()
	})
	i, s := f()
	assert.Equal(t, i, 0)
	assert.Equal(t, s, "")
}
