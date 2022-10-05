package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	hashmultierr "github.com/hashicorp/go-multierror"
	"github.com/ovsinc/errors"
	"github.com/stretchr/testify/require"
	ubermulierr "go.uber.org/multierr"
	"golang.org/x/xerrors"
)

//
// stderrors
//

// Vendors

func BenchmarkVendorStandartError(b *testing.B) {
	e := stderrors.New("hello1")

	require.Equal(b, e.Error(), "hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkVendorStandartConstructor(b *testing.B) {
	e := stderrors.New("hello1")

	require.Equal(b, e.Error(), "hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stderrors.New("hello1")
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

func BenchmarkVendorXerrorsConstructor(b *testing.B) {
	e := xerrors.New("hello1")

	require.Equal(b, e.Error(), "hello1")

	for i := 0; i < b.N; i++ {
		_ = xerrors.New("hello1")
	}
}

// MY

func BenchmarkVendorMyNewNormal(b *testing.B) {
	err := errors.New("hello1")

	require.Equal(b, err.Error(), "hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorMyConstructorNormal(b *testing.B) {
	err := errors.New("hello1")

	require.Equal(b, err.Error(), "hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.New("hello1")
	}
}

func BenchmarkVendorMyConstructorFullNoCtx(b *testing.B) {
	err := errors.NewWith(
		errors.SetMsg("hello1"),
		errors.SetID("IDhello1"),
		errors.SetOperation("nothing"),
		errors.SetErrorType("myerrtype"),
	)

	require.Equal(b, err.Error(), "(myerrtype) [nothing] hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.NewWith(
			errors.SetMsg("hello1"),
			errors.SetID("IDhello1"),
			errors.SetOperation("nothing"),
			errors.SetErrorType("myerrtype"),
		)
	}
}

func BenchmarkVendorMyNewFullNoCtx(b *testing.B) {
	err := errors.NewWith(
		errors.SetMsg("hello1"),
		errors.SetID("IDhello1"),
		errors.SetOperation("nothing"),
		errors.SetErrorType("myerrtype"),
	)

	require.Equal(b, err.Error(), "(myerrtype) [nothing] hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorMyNewFull(b *testing.B) {
	err := errors.NewWith(
		errors.SetMsg("hello1"),
		errors.AppendContextInfo("hello", "world"),
		errors.SetID("IDhello1"),
		errors.SetOperation("nothing"),
	)

	require.Equal(b, err.Error(), "[nothing] {hello:world} hello1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

//
// multierr
//

// MY

func BenchmarkVendorMyMulti2StdErr(b *testing.B) {
	err := errors.Wrap(
		stderrors.New("hello1"),
		stderrors.New("hello2"),
	)

	require.Equal(b, err.Error(), "the following errors occurred:\n\t#1 hello1\n\t#2 hello2\n")

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorMyMulti2StdErrConstructor(b *testing.B) {
	err := errors.Wrap(
		stderrors.New("hello1"),
		stderrors.New("hello2"),
	)

	require.Equal(b, err.Error(), "the following errors occurred:\n\t#1 hello1\n\t#2 hello2\n")

	for i := 0; i < b.N; i++ {
		_ = errors.Wrap(
			stderrors.New("hello1"),
			stderrors.New("hello2"),
		)
	}
}

func BenchmarkVendorMyMulti2MySimple(b *testing.B) {
	err := errors.Wrap(
		errors.New("hello1"),
		errors.New("hello2"),
	)

	require.Equal(b, err.Error(), "the following errors occurred:\n\t#1 hello1\n\t#2 hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkVendorMyMulti2MySimpleConstructor(b *testing.B) {
	err := errors.Wrap(
		errors.New("hello1"),
		errors.New("hello2"),
	)

	require.Equal(b, err.Error(), "the following errors occurred:\n\t#1 hello1\n\t#2 hello2\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.Wrap(
			errors.New("hello1"),
			errors.New("hello2"),
		)
	}
}

// Vendor

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

func BenchmarkVendorHashiMulti2StdErrConstructor(b *testing.B) {
	err := hashmultierr.Append(
		stderrors.New("hello1"),
		stderrors.New("hello2"),
	)

	require.Equal(b, err.Error(), "2 errors occurred:\n\t* hello1\n\t* hello2\n\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hashmultierr.Append(
			stderrors.New("hello1"),
			stderrors.New("hello2"),
		)
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

func BenchmarkVendorUberMulti2StdErrConstructor(b *testing.B) {
	err := ubermulierr.Append(
		stderrors.New("hello1"),
		stderrors.New("hello2"),
	)

	require.Equal(b, err.Error(), "hello1; hello2")

	for i := 0; i < b.N; i++ {
		_ = ubermulierr.Append(
			stderrors.New("hello1"),
			stderrors.New("hello2"),
		)
	}
}

func BenchmarkVendorUberMulti2MyErr(b *testing.B) {
	err := ubermulierr.Append(
		errors.New("hello1"),
		errors.New("hello2"),
	)

	require.Equal(b, err.Error(), "hello1; hello2")

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
