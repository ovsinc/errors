package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"gitlab.com/ovsinc/errors/log"
)

// defaultFormatFn функция форматирования, используемая по-умолчанию
var defaultFormatFn FormatFn = StringFormat //nolint:gochecknoglobals

// FormatFn тип функции форматирования.
type FormatFn func(buf io.Writer, e *Error)

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
		op0 := ops[0].String()
		_, _ = io.WriteString(buf, "\"")
		_, _ = io.WriteString(buf, op0)
		_, _ = io.WriteString(buf, "\"")
		for _, s := range ops[1:] {
			_, _ = io.WriteString(buf, ",")
			opN := s.String()
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

	if e.ErrorType() != "" {
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
		op0 := ops[0].String()
		_, _ = io.WriteString(buf, op0)
		for _, s := range ops[1:] {
			_, _ = io.WriteString(buf, ",")
			opN := s.String()
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
		_, _ = io.WriteString(buf, " -- ")
	}

	_ = e.writeTranslate(buf, msg)
}

//

// Format производит форматирование строки, для поддержки fmt.Printf().
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'c':
		fmt.Fprintf(s, "%v\n", e.ContextInfo())

	case 'o':
		fmt.Fprintf(s, "%v\n", e.Operations())

	case 'l':
		_, _ = io.WriteString(s, e.Severity().String())

	case 't':
		_, _ = io.WriteString(s, e.ErrorType())

	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.Sdump())
			return
		}
		_, _ = io.WriteString(s, e.Error())

	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}
