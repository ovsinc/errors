package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/ovsinc/errors"
	"gitlab.com/ovsinc/errors/log"

	hashmultierr "github.com/hashicorp/go-multierror"

	ubermulierr "go.uber.org/multierr"
)

func BenchmarkStandartError(b *testing.B) {
	e := stderrors.New("[UNKNOWN_TYPE][ERROR] -- hello1")

	require.Equal(b, e.Error(), "[UNKNOWN_TYPE][ERROR] -- hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkFmt(b *testing.B) {
	e := fmt.Errorf("[UNKNOWN_TYPE][ERROR] -- hello1")

	require.Equal(b, e.Error(), "[UNKNOWN_TYPE][ERROR] -- hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkMyNewMsgOnly(b *testing.B) {
	err := errors.New(
		"hello1",
		errors.SetErrorType(""),
		errors.SetSeverity(log.SeverityUnknown),
	)

	require.Equal(b, err.Error(), "hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkMyNew(b *testing.B) {
	err := errors.New(
		"hello1",
	)

	require.Equal(b, err.Error(), "[UNKNOWN_TYPE][ERROR] -- hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

// mutierr

func BenchmarkMyMulti2Err(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := errors.Append(
		stderrors.New("[UNKNOWN_TYPE][ERROR] -- hello1"),
		stderrors.New("[UNKNOWN_TYPE][ERROR] -- hello2"),
	)

	require.Equal(b, err.Error(), "* [UNKNOWN_TYPE][ERROR] -- hello1\n* [UNKNOWN_TYPE][ERROR] -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkMyMultiMsgOnly2Err(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := errors.Append(
		errors.New(
			"hello1",
			errors.SetErrorType(""),
			errors.SetSeverity(log.SeverityUnknown),
		),
		errors.New(
			"hello2",
			errors.SetErrorType(""),
			errors.SetSeverity(log.SeverityUnknown),
		),
	)

	require.Equal(b, err.Error(), "* hello1\n* hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkHashiMulti2Err(b *testing.B) {
	err := hashmultierr.Append(
		stderrors.New("[UNKNOWN_TYPE][ERROR] -- hello1"),
		stderrors.New("[UNKNOWN_TYPE][ERROR] -- hello2"),
	)

	require.Equal(b, err.Error(), "2 errors occurred:\n\t* [UNKNOWN_TYPE][ERROR] -- hello1\n\t* [UNKNOWN_TYPE][ERROR] -- hello2\n\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkUberMulti2Err(b *testing.B) {
	err := ubermulierr.Append(
		stderrors.New("[UNKNOWN_TYPE][ERROR] -- hello1"),
		stderrors.New("[UNKNOWN_TYPE][ERROR] -- hello2"),
	)

	require.Equal(b, err.Error(), "[UNKNOWN_TYPE][ERROR] -- hello1; [UNKNOWN_TYPE][ERROR] -- hello2")

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
