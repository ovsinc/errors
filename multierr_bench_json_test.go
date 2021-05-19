package errors_test

import (
	"testing"

	"github.com/ovsinc/errors"
	"github.com/stretchr/testify/require"
)

var (
	je1 = errors.New(
		"hello1",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	je2 = errors.New(
		"hello2",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello2": "world", "my2": "name"}),
	)

	je3 = errors.New(
		"hello3",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello3": "world", "my3": "name"}),
	)
)

func BenchmarkJsonFn(b *testing.B) {
	errors.DefaultFormatFn = errors.JSONFormat

	e := errors.New(
		"hello",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "hi": "there"}),
		errors.SetFormatFn(errors.JSONFormat),
	)

	require.JSONEq(b, e.Error(), "{\"id\":\"\", \"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"hi\":\"there\"},\"msg\":\"hello\"}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat3Errs(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat
	errors.DefaultFormatFn = errors.JSONFormat

	e := errors.Append(je1, je2, je3)

	require.JSONEq(b, e.Error(), "{\"count\":3,\"messages\":[{\"id\":\"\", \"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"my\":\"name\"},\"msg\":\"hello1\"},{\"id\":\"\", \"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"},{\"id\":\"\", \"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello3\":\"world\",\"my3\":\"name\"},\"msg\":\"hello3\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat2Errs(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat
	errors.DefaultFormatFn = errors.JSONFormat

	e := errors.Wrap(je1, je2)

	require.JSONEq(b, e.Error(), "{\"count\":2,\"messages\":[{\"id\":\"\", \"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"my\":\"name\"},\"msg\":\"hello1\"},{\"id\":\"\", \"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat1Err(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat
	errors.DefaultFormatFn = errors.JSONFormat

	e := errors.Wrap(nil, je2)

	require.JSONEq(b, e.Error(), "{\"id\":\"\", \"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}
