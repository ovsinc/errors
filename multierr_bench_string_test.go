package errors_test

import (
	"fmt"
	"testing"

	"github.com/ovsinc/errors"
	"github.com/stretchr/testify/require"
)

var (
	se1 = errors.New(
		"hello1",
		errors.SetOperation("write"),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	se2 = errors.New(
		"hello2",
		errors.SetOperation("read"),
		errors.SetContextInfo(errors.CtxMap{"hello2": "world", "my2": "name"}),
	)

	se3 = errors.New(
		"hello3",
		errors.SetOperation("read"),
		errors.SetContextInfo(errors.CtxMap{"hello3": "world", "my3": "name"}),
	)
)

func BenchmarkStringFn(b *testing.B) {
	e := errors.New(
		"hello",
		errors.SetOperation("write"),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	data, err := e.Marshal(&errors.MarshalString{})
	require.Nil(b, err)

	require.Equal(b, string(data), "write: {hello:world,my:name} -- hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkFormatFmt(b *testing.B) {
	e := errors.New(
		"hello",
		errors.SetOperation("write"),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "name": "john"}),
	)

	data, err := e.Marshal(&errors.MarshalString{})
	require.Nil(b, err)

	require.Equal(b, string(data), "write: {hello:world,name:john} -- hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// fmt.Fprintf(ioutil.Discard, "%v", e)
		_ = fmt.Sprintf("%v", e)
	}
}

func BenchmarkStringMultierrFuncFormat3Errs(b *testing.B) {
	e := errors.Combine(se1, se2, se3)

	marshal := &errors.MarshalString{}
	data, err := marshal.Marshal(e)
	require.Nil(b, err)

	require.Equal(b, string(data), "the following errors occurred:\n\t#1 write: {hello:world,my:name} -- hello1\n\t#2 read: {hello2:world,my2:name} -- hello2\n\t#3 read: {hello3:world,my3:name} -- hello3\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkStringMultierrFuncFormat2Errs(b *testing.B) {
	e := errors.Wrap(se1, se2)

	marshal := &errors.MarshalString{}
	data, err := marshal.Marshal(e)
	require.Nil(b, err)

	require.Equal(b, string(data), "the following errors occurred:\n\t#1 write: {hello:world,my:name} -- hello1\n\t#2 read: {hello2:world,my2:name} -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkStringMultierrFuncFormat1Err(b *testing.B) {
	e := errors.Wrap(nil, se2)

	marshal := &errors.MarshalString{}
	data, err := marshal.Marshal(e)
	require.Nil(b, err)

	require.Equal(b, string(data), "the following errors occurred:\n\t#1 read: {hello2:world,my2:name} -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}
