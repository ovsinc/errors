package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"

	"gitlab.com/ovsinc/errors/log"
)

var (
	// defaultFormatFn функция форматирования, используемая по-умолчанию
	defaultFormatFn FormatFn //nolint:gochecknoglobals

	// DefaultMultierrFormatFunc функция форматирования для multierr ошибок.
	DefaultMultierrFormatFunc MultierrFormatFn //nolint:gochecknoglobals

	_multilinePrefix    = []byte("the following errors occurred:")
	_multilineSeparator = []byte("\n")
	_multilineIndent    = []byte("* ")
	_msgSeparator       = []byte(" -- ")
)

type (
	// FormatFn тип функции форматирования.
	FormatFn func(w io.Writer, e *Error)

	// MultierrFormatFn типу функции морматирования для multierr.
	MultierrFormatFn func(w io.Writer, es []error)
)

// JSONFormat функция форматирования вывода сообщения *Error в JSON.
func JSONFormat(buf io.Writer, e *Error) {
	if e == nil {
		_, _ = io.WriteString(buf, "null")
		return
	}

	_, _ = io.WriteString(buf, "{")

	// ErrorType
	_, _ = io.WriteString(buf, "\"error_type\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = io.WriteString(buf, e.errorType)
	_, _ = io.WriteString(buf, "\",")

	// Severity
	_, _ = io.WriteString(buf, "\"severity\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = io.WriteString(buf, e.severity.String())
	_, _ = io.WriteString(buf, "\",")

	// Operations
	_, _ = io.WriteString(buf, "\"operations\":[")
	ops := e.Operations()
	if len(ops) > 0 {
		op0 := ops[0]
		_, _ = io.WriteString(buf, "\"")
		_, _ = io.WriteString(buf, op0)
		_, _ = io.WriteString(buf, "\"")
		for _, opN := range ops[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = io.WriteString(buf, "\"")
			_, _ = io.WriteString(buf, opN)
			_, _ = io.WriteString(buf, "\"")
		}
	}
	_, _ = io.WriteString(buf, "],")

	// ContextInfo
	_, _ = io.WriteString(buf, "\"context\":")
	data, _ := json.Marshal(e.ContextInfo())
	_, _ = buf.Write(data)
	_, _ = io.WriteString(buf, ",")

	// Msg
	_, _ = io.WriteString(buf, "\"msg\":")
	_, _ = io.WriteString(buf, "\"")
	if len(e.msg) == 0 {
		_, _ = io.WriteString(buf, "null")
	} else {
		_ = e.writeTranslate(buf, e.msg)
	}
	_, _ = io.WriteString(buf, "\"")

	_, _ = io.WriteString(buf, "}")
}

// StringFormat функция форматирования вывода сообщения *Error в виде строки.
// Используется по-умолчанию.
func StringFormat(buf io.Writer, e *Error) { //nolint:cyclop
	if e == nil {
		return
	}

	writeDelim := false

	ops := e.Operations()
	ctxs := e.ContextInfo()
	msg := e.Msg()

	if e.errorType != "" {
		_, _ = io.WriteString(buf, "[")
		_, _ = io.WriteString(buf, e.errorType)
		_, _ = io.WriteString(buf, "]")
		writeDelim = true
	}

	if e.Severity() > log.SeverityUnknown {
		_, _ = io.WriteString(buf, "[")
		_, _ = io.WriteString(buf, e.severity.String())
		_, _ = io.WriteString(buf, "]")
		writeDelim = true
	}

	if len(ops) > 0 {
		_, _ = io.WriteString(buf, "[")
		op0 := ops[0]
		_, _ = io.WriteString(buf, op0)
		for _, opN := range ops[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = io.WriteString(buf, opN)
		}
		_, _ = io.WriteString(buf, "]")
		writeDelim = true
	}

	if len(ctxs) > 0 {
		_, _ = io.WriteString(buf, "<")
		ctxskeys := make([]string, 0, len(ctxs))
		for i := range ctxs {
			ctxskeys = append(ctxskeys, i)
		}
		sort.Strings(ctxskeys)
		_, _ = fmt.Fprintf(buf, "%s:%v", ctxskeys[0], ctxs[ctxskeys[0]])
		for _, i := range ctxskeys[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = fmt.Fprintf(buf, "%s:%v", i, ctxs[i])
		}
		_, _ = io.WriteString(buf, ">")
		writeDelim = true
	}

	if writeDelim && len(msg) > 0 {
		_, _ = buf.Write(_msgSeparator)
	}

	_ = e.writeTranslate(buf, msg)
}

//

// multierr

// StringMultierrFormatFunc функция форматирования вывода сообщения для multierr в виде строки.
// Используется по-умолчанию.
func StringMultierrFormatFunc(w io.Writer, es []error) {
	if len(es) == 0 {
		_, _ = io.WriteString(w, "")
		return
	}

	_, _ = w.Write(_multilinePrefix)
	_, _ = w.Write(_multilineSeparator)

	for _, err := range es {
		_, _ = w.Write(_multilineIndent)
		_, _ = io.WriteString(w, err.Error())
		_, _ = w.Write(_multilineSeparator)
	}
}

// JSONMultierrFuncFormat функция форматирования вывода сообщения для multierr в виде JSON.
func JSONMultierrFuncFormat(w io.Writer, es []error) {
	if len(es) == 0 {
		_, _ = io.WriteString(w, "null")
	}

	_, _ = io.WriteString(w, "{")

	_, _ = io.WriteString(w, "\"count\":")
	_, _ = io.WriteString(w, strconv.Itoa(len(es)))
	_, _ = io.WriteString(w, ",")

	_, _ = io.WriteString(w, "\"messages\":")
	_, _ = io.WriteString(w, "[")
	writeErrFn := func(e error) {
		var myerr *Error
		if As(e, &myerr) {
			JSONFormat(w, myerr)
		} else {
			_, _ = fmt.Fprintf(w, "\"%v\"", myerr)
		}
	}
	switch len(es) {
	case 0:
	case 1:
		writeErrFn(es[0])
	default:
		writeErrFn(es[0])
		for _, e := range es[1:] {
			_, _ = io.WriteString(w, ",")
			writeErrFn(e)
		}
	}
	_, _ = io.WriteString(w, "]")

	_, _ = io.WriteString(w, "}")
}
