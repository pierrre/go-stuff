package main

import (
	"fmt"
	"io"
	"iter"
	"math/big"
	"os"

	"github.com/pierrre/errors"
	"github.com/pierrre/go-libs/bufpool"
)

func main() {
	err := FibonacciWrite(os.Stdout, 1000)
	if err != nil {
		panic(err)
	}
}

func Fibonacci(n int) *big.Int {
	var i int
	var v *big.Int
	for i, v = range FibonacciSeq() {
		if i >= n {
			break
		}
	}
	return v
}

func FibonacciString(n int) string {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	_ = FibonacciWrite(buf, n) // Ignore error because it's a bytes.Buffer.
	return buf.String()
}

var bufPool = &bufpool.Pool{}

func FibonacciSeq() iter.Seq2[int, *big.Int] {
	return func(yield func(int, *big.Int) bool) {
		a, b := big.NewInt(0), big.NewInt(1)
		for i := 0; ; i++ {
			if !yield(i, a) {
				return
			}
			a.Add(a, b)
			a, b = b, a
		}
	}
}

func FibonacciWrite(w io.Writer, n int) error {
	for i, v := range FibonacciSeq() {
		if i >= n {
			break
		}
		// TODO reduce allocations in math/big: https://github.com/golang/go/issues/71465
		_, err := fmt.Fprintf(w, "%v\n", v)
		if err != nil {
			return errors.Wrapf(err, "write %d iteration", i)
		}
	}
	return nil
}
