package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/ovsinc/errors"
	"gitlab.com/ovsinc/errors/log"
)

func BenchmarkStringFn(b *testing.B) {
	e := errors.New(
		"hello",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(log.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	require.Equal(b, e.Error(), "[not found][ERROR][write]<hello:world,my:name> -- hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkFormatFmt(b *testing.B) {
	e := errors.New(
		"hello",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(log.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "name": "john"}),
	)

	require.Equal(b, e.Error(), "[not found][ERROR][write]<hello:world,name:john> -- hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// fmt.Fprintf(ioutil.Discard, "%v", e)
		_ = fmt.Sprintf("%v", e)
	}
}

func BenchmarkStringMultierrFuncFormat3Errs(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	e := errors.Append(e1, e2, e3)

	require.Equal(b, e.Error(), "the following errors occurred:\n* [not found][ERROR][write]<hello:world,my:name> -- hello1\n* [not found][ERROR][read]<hello2:world,my2:name> -- hello2\n* [not found][ERROR][read]<hello3:world,my3:name> -- hello3\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkStringMultierrFuncFormat2Errs(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	e := errors.Wrap(e1, e2)

	require.Equal(b, e.Error(), "the following errors occurred:\n* [not found][ERROR][write]<hello:world,my:name> -- hello1\n* [not found][ERROR][read]<hello2:world,my2:name> -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkStringMultierrFuncFormat1Err(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	e := errors.Wrap(nil, e2)

	require.Equal(b, e.Error(), "the following errors occurred:\n* [not found][ERROR][read]<hello2:world,my2:name> -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}
