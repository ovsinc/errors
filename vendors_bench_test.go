// +build vendors

package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	hashmultierr "github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
	"gitlab.com/ovsinc/errors"
	ubermulierr "go.uber.org/multierr"
	"golang.org/x/xerrors"
)

func BenchmarkVendorStandartError(b *testing.B) {
	e := stderrors.New("hello1")

	require.Equal(b, e.Error(), "hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkVendorFmt(b *testing.B) {
	e := fmt.Errorf("%s", "hello1")

	require.Equal(b, e.Error(), "hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkVendorXerrors(b *testing.B) {
	e := xerrors.New("hello1")

	require.Equal(b, e.Error(), "hello1")

	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkVendorMyNewNormal(b *testing.B) {
	err := errors.New(
		"hello1",
	)

	require.Equal(b, err.Error(), "hello1")

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
		errors.SetErrorType("not found"),
	)

	require.Equal(b, err.Error(), "(not found)[nothing]<hello:world> -- hello1")

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

	require.Equal(b, err.Error(), "[nothing]<hello:world> -- У John Snow имеется 5 непрочитанных сообщений.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

// mutierr

func BenchmarkVendorMyMulti2StdErr(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := errors.Wrap(
		stderrors.New("hello1"),
		stderrors.New("hello2"),
	)

	require.Equal(b, err.Error(), "the following errors occurred:\n* hello1\n* hello2\n")

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

	require.Equal(b, err.Error(), "the following errors occurred:\n* hello1\n* hello2\n")

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
			errors.SetSeverity(errors.SeverityUnknown),
		),
		errors.New(
			"hello2",
			errors.SetErrorType(""),
			errors.SetSeverity(errors.SeverityUnknown),
		),
	)

	require.Equal(b, err.Error(), "the following errors occurred:\n* hello1\n* hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorHashiMulti2StdErr(b *testing.B) {
	err := hashmultierr.Append(
		stderrors.New("hello1"),
		stderrors.New("hello2"),
	)

	require.Equal(b, err.Error(), "2 errors occurred:\n\t* hello1\n\t* hello2\n\n")

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

	require.Equal(b, err.Error(), "2 errors occurred:\n\t* hello1\n\t* hello2\n\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorUberMulti2StdErr(b *testing.B) {
	err := ubermulierr.Append(
		stderrors.New("hello1"),
		stderrors.New("hello2"),
	)

	require.Equal(b, err.Error(), "hello1; hello2")

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorUberMulti2MyNormalErr(b *testing.B) {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := ubermulierr.Append(
		errors.New("hello1"),
		errors.New("hello2"),
	)

	require.Equal(b, err.Error(), "hello1; hello2")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
