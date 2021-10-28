package errors

import (
	"io"
)

var (
	// DefaultFormatFn функция форматирования, используемая по-умолчанию
	DefaultFormatFn FormatFn //nolint:gochecknoglobals

	// DefaultMultierrFormatFunc функция форматирования для multierr ошибок.
	DefaultMultierrFormatFunc MultierrFormatFn //nolint:gochecknoglobals
)

type (
	// FormatFn тип функции форматирования.
	FormatFn func(w io.Writer, e *Error)

	// MultierrFormatFn типу функции морматирования для multierr.
	MultierrFormatFn func(w io.Writer, es []error)
)
