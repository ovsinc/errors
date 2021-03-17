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

func BenchmarkVendorStandartError(b *testing.B) {
	e := stderrors.New("[ERROR] -- hello1")

	require.Equal(b, e.Error(), "[ERROR] -- hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkVendorFmt(b *testing.B) {
	e := fmt.Errorf("[ERROR] -- hello1")

	require.Equal(b, e.Error(), "[ERROR] -- hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkVendorMyNewMsgOnly(b *testing.B) {
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

func BenchmarkVendorMyNewNormal(b *testing.B) {
	err := errors.New(
		"hello1",
	)

	require.Equal(b, err.Error(), "[ERROR] -- hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorMyNewFull(b *testing.B) {
	err := errors.New(
		"hello1",
		errors.AppendContextInfo("hello", "world"),
		errors.SetID("IDhello1"),
		errors.SetOperations("nothing"),
	)

	require.Equal(b, err.Error(), "[ERROR][nothing]<hello:world> -- hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorMyNewWithTranslate(b *testing.B) {
	errEmailsUnreadMsg := localTransContext()
	localizer := localizePrepare()

	err := errors.New(
		"hello1",
		errors.AppendContextInfo("hello", "world"),
		errors.SetOperations("nothing"),
		errors.SetID("ErrEmailsUnreadMsg"),
		errors.SetTranslateContext(&errEmailsUnreadMsg),
		errors.SetLocalizer(localizer),
	)

	require.Equal(b, err.Error(), "[ERROR][nothing]<hello:world> -- У John Snow имеется 5 непрочитанных сообщений.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

// mutierr

func BenchmarkVendorMyMulti2StdErr(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := errors.Wrap(
		stderrors.New("[ERROR] -- hello1"),
		stderrors.New("[ERROR] -- hello2"),
	)

	require.Equal(b, err.Error(), "* [ERROR] -- hello1\n* [ERROR] -- hello2\n")

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorMyMulti2ErrNormal(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := errors.Wrap(
		errors.New("hello1"),
		errors.New("hello2"),
	)

	require.Equal(b, err.Error(), "* [ERROR] -- hello1\n* [ERROR] -- hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorMyMulti2ErrMsgOnly(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := errors.Wrap(
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

func BenchmarkVendorHashiMulti2StdErr(b *testing.B) {
	err := hashmultierr.Append(
		stderrors.New("[ERROR] -- hello1"),
		stderrors.New("[ERROR] -- hello2"),
	)

	require.Equal(b, err.Error(), "2 errors occurred:\n\t* [ERROR] -- hello1\n\t* [ERROR] -- hello2\n\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorHashiMulti2MyErr(b *testing.B) {
	err := hashmultierr.Append(
		errors.New("hello1"),
		errors.New("hello2"),
	)

	require.Equal(b, err.Error(), "2 errors occurred:\n\t* [ERROR] -- hello1\n\t* [ERROR] -- hello2\n\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorUberMulti2StdErr(b *testing.B) {
	err := ubermulierr.Append(
		stderrors.New("[ERROR] -- hello1"),
		stderrors.New("[ERROR] -- hello2"),
	)

	require.Equal(b, err.Error(), "[ERROR] -- hello1; [ERROR] -- hello2")

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorUberMulti2MyErr(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := ubermulierr.Append(
		errors.New("hello1"),
		errors.New("hello2"),
	)

	require.Equal(b, err.Error(), "[ERROR] -- hello1; [ERROR] -- hello2")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
