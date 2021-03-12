package errors

import (
	"testing"

	"gitlab.com/ovsinc/errors/log"

	"github.com/stretchr/testify/require"
)

var (
	e1 = New(
		"hello1",
		SetErrorType(NewErrorType("not found")),
		SetOperations(NewOperation("write")),
		SetSeverity(log.SeverityError),
		SetContextInfo(CtxMap{"hello": "world", "my": "name"}),
	)

	e2 = New(
		"hello2",
		SetErrorType(NewErrorType("not found")),
		SetOperations(NewOperation("read")),
		SetSeverity(log.SeverityError),
		SetContextInfo(CtxMap{"hello2": "world", "my2": "name"}),
	)

	e3 = New(
		"hello3",
		SetErrorType(NewErrorType("not found")),
		SetOperations(NewOperation("read")),
		SetSeverity(log.SeverityError),
		SetContextInfo(CtxMap{"hello3": "world", "my3": "name"}),
	)
)

func BenchmarkStringMultierrFormatFunc3Errs(b *testing.B) {
	DefaultMultierrFormatFunc = StringMultierrFormatFunc

	e := Append(e1, e2, e3)

	require.Equal(b, e.Error(), "* [not found][ERROR][write]<hello:world,my:name> -- hello1\n* [not found][ERROR][read]<hello2:world,my2:name> -- hello2\n* [not found][ERROR][read]<hello3:world,my3:name> -- hello3\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkStringMultierrFormatFunc2Errs(b *testing.B) {
	DefaultMultierrFormatFunc = StringMultierrFormatFunc

	e := Wrap(e1, e2)

	require.Equal(b, e.Error(), "* [not found][ERROR][write]<hello:world,my:name> -- hello1\n* [not found][ERROR][read]<hello2:world,my2:name> -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkStringMultierrFormatFunc1Err(b *testing.B) {
	DefaultMultierrFormatFunc = StringMultierrFormatFunc

	e := Wrap(nil, e2)

	require.Equal(b, e.Error(), "* [not found][ERROR][read]<hello2:world,my2:name> -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat3Errs(b *testing.B) {
	DefaultMultierrFormatFunc = JSONMultierrFuncFormat

	e := Append(e1, e2, e3)

	require.JSONEq(b, e.Error(), "{\"count\":3,\"messages\":[{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"my\":\"name\"},\"msg\":\"hello1\"},{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"},{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello3\":\"world\",\"my3\":\"name\"},\"msg\":\"hello3\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat2Errs(b *testing.B) {
	DefaultMultierrFormatFunc = JSONMultierrFuncFormat

	e := Wrap(e1, e2)

	require.JSONEq(b, e.Error(), "{\"count\":2,\"messages\":[{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"my\":\"name\"},\"msg\":\"hello1\"},{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat1Err(b *testing.B) {
	DefaultMultierrFormatFunc = JSONMultierrFuncFormat

	e := Wrap(nil, e2)

	require.JSONEq(b, e.Error(), "{\"count\":1,\"messages\":[{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}
