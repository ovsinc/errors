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
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	se2 = errors.New(
		"hello2",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello2": "world", "my2": "name"}),
	)

	se3 = errors.New(
		"hello3",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello3": "world", "my3": "name"}),
	)
)

func BenchmarkStringFn(b *testing.B) {
	errors.DefaultFormatFn = errors.StringFormat

	e := errors.New(
		"hello",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	require.Equal(b, e.Error(), "(not found)[write]<hello:world,my:name> -- hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkFormatFmt(b *testing.B) {
	errors.DefaultFormatFn = errors.StringFormat

	e := errors.New(
		"hello",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "name": "john"}),
	)

	require.Equal(b, e.Error(), "(not found)[write]<hello:world,name:john> -- hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// fmt.Fprintf(ioutil.Discard, "%v", e)
		_ = fmt.Sprintf("%v", e)
	}
}

func BenchmarkStringMultierrFuncFormat3Errs(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc
	errors.DefaultFormatFn = errors.StringFormat

	e := errors.Append(se1, se2, se3)

	require.Equal(b, e.Error(), "the following errors occurred:\n* (not found)[write]<hello:world,my:name> -- hello1\n* (not found)[read]<hello2:world,my2:name> -- hello2\n* (not found)[read]<hello3:world,my3:name> -- hello3\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkStringMultierrFuncFormat2Errs(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc
	errors.DefaultFormatFn = errors.StringFormat

	e := errors.Wrap(se1, se2)

	require.Equal(b, e.Error(), "the following errors occurred:\n* (not found)[write]<hello:world,my:name> -- hello1\n* (not found)[read]<hello2:world,my2:name> -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkStringMultierrFuncFormat1Err(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc
	errors.DefaultFormatFn = errors.StringFormat

	e := errors.Wrap(nil, se2)

	require.Equal(b, e.Error(), "(not found)[read]<hello2:world,my2:name> -- hello2")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}
