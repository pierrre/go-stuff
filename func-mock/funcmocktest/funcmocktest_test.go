package funcmocktest

import (
	"testing"

	"github.com/pierrre/assert"
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
	}, assert.Report(asserttest.NewReportAuto(t)))
	f()
	f()
	f()
}

func TestMockListFailRemaningCalls(t *testing.T) {
	f := MockList(t, []func(){
		func() {},
		func() {},
		func() {},
	}, assert.Report(asserttest.NewReportAuto(t)))
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
	f := MockCount(t, func() {}, 2, assert.Report(asserttest.NewReportAuto(t)))
	f()
	f()
	f()
}
