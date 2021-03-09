package errors

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/ovsinc/errors/log"
)

func BenchmarkStringFn(b *testing.B) {
	e := New(
		"hello",
		SetErrorType(NewErrorType("not found")),
		SetOperations(NewOperation("write")),
		SetSeverity(log.SeverityError),
		SetContextInfo(CtxMap{"hello": "world", "my": "name"}),
	)

	require.Equal(b, e.Error(), "[not found][ERROR][write]<hello:world,my:name> -- hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkFormatFmt(b *testing.B) {
	e := New(
		"hello",
		SetErrorType(NewErrorType("not found")),
		SetOperations(NewOperation("write")),
		SetSeverity(log.SeverityError),
		SetContextInfo(CtxMap{"hello": "world", "name": "john"}),
	)

	require.Equal(b, e.Error(), "[not found][ERROR][write]<hello:world,name:john> -- hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(ioutil.Discard, "%v", e)
	}
}

func BenchmarkJsonFn(b *testing.B) {
	e := New(
		"hello",
		SetErrorType(NewErrorType("not found")),
		SetOperations(NewOperation("write")),
		SetSeverity(log.SeverityError),
		SetContextInfo(CtxMap{"hello": "world", "hi": "there"}),
		SetFormatFn(JSONFormat),
	)

	require.JSONEq(b, e.Error(), "{\"error_type\":\"not found\",\"severity\":\"ERROR\",\"operations\":[\"write\"],\"context\":{\"hello\":\"world\",\"hi\":\"there\"},\"msg\":\"hello\"}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}
