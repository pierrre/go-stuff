package funcmocktest

import (
	"testing"

	"github.com/pierrre/assert/asserttest"
)

func TestMockList(t *testing.T) {
	f := MockList(t, []func(){
		func() {},
		func() {},
		func() {},
	})
	f()
	f()
	f()
}

func TestMockListFailTooManyCalls(t *testing.T) {
	f := MockList(t, []func(){
		func() {},
		func() {},
	}, asserttest.ReportAuto(t))
	f()
	f()
	f()
}

func TestMockListFailRemaningCalls(t *testing.T) {
	f := MockList(t, []func(){
		func() {},
		func() {},
		func() {},
	}, asserttest.ReportAuto(t))
	f()
	f()
}

func TestMockCount(t *testing.T) {
	f := MockCount(t, func() {}, 3)
	f()
	f()
	f()
}

func TestMockCountFail(t *testing.T) {
	f := MockCount(t, func() {}, 2, asserttest.ReportAuto(t))
	f()
	f()
	f()
}
