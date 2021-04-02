package errors_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/ovsinc/errors"
	"gitlab.com/ovsinc/errors/log"
)

var (
	e1 = errors.New(
		"hello1",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(log.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	e2 = errors.New(
		"hello2",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(log.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello2": "world", "my2": "name"}),
	)

	e3 = errors.New(
		"hello3",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(log.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello3": "world", "my3": "name"}),
	)
)

func BenchmarkJsonFn(b *testing.B) {
	e := errors.New(
		"hello",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(log.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "hi": "there"}),
		errors.SetFormatFn(errors.JSONFormat),
	)

	require.JSONEq(b, e.Error(), "{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"hi\":\"there\"},\"msg\":\"hello\"}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat3Errs(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat

	e := errors.Append(e1, e2, e3)

	require.JSONEq(b, e.Error(), "{\"count\":3,\"messages\":[{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"my\":\"name\"},\"msg\":\"hello1\"},{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"},{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello3\":\"world\",\"my3\":\"name\"},\"msg\":\"hello3\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat2Errs(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat

	e := errors.Wrap(e1, e2)

	require.JSONEq(b, e.Error(), "{\"count\":2,\"messages\":[{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"my\":\"name\"},\"msg\":\"hello1\"},{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat1Err(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat

	e := errors.Wrap(nil, e2)

	require.JSONEq(b, e.Error(), "{\"count\":1,\"messages\":[{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"read\"],\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}
