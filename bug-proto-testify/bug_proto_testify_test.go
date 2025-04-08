package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/proto"
)

func Test(t *testing.T) {
	a := &money.Money{
		CurrencyCode: "EUR",
		Units:        123,
		Nanos:        456000000,
	}
	b := &money.Money{
		CurrencyCode: "EUR",
		Units:        123,
		Nanos:        456000000,
	}
	require.Equal(t, a, b)
	_ = proto.CloneOf(b)
	//require.Equal(t, a, b)
}
