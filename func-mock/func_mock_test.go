package funcmock

import (
	"testing"
)

func TestMock(t *testing.T) {
	f := Mock([]func(s string) int{
		func(s string) int {
			t.Log("Hello", s)
			return 1
		},
		func(s string) int {
			t.Log("World", s)
			return 2
		},
	})
	t.Log(f("aaaaa"))
	t.Log(f("bbbbb"))
	t.Log(f("ccccc"))
}
