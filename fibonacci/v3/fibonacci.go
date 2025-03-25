//nolint:forbidigo // Print text.
package main

import (
	"fmt"
	"iter"
	"math/big"
)

func main() {
	s := fibonacciString(1000)
	fmt.Println(s)
}

func fibonacci(n int) *big.Int {
	var i int
	var v *big.Int
	for i, v = range fibonacciSeq() {
		if i >= n {
			break
		}
	}
	return v
}

func fibonacciString(n int) string {
	s := ""
	for i, v := range fibonacciSeq() {
		if i >= n {
			break
		}
		vs := v.String()
		s += vs + "\n"
	}
	return s
}

func fibonacciSeq() iter.Seq2[int, *big.Int] {
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
