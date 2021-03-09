package errors

import (
	"encoding/json"
	"fmt"
	"sort"

	"io"

	"github.com/valyala/bytebufferpool"
	"gitlab.com/ovsinc/errors/log"
)

// defaultFormatFn функция форматирования, используемая по-умолчанию
var defaultFormatFn FormatFn = StringFormat

// FormatFn функция форматирования
type FormatFn func(e *Error) string

func JSONFormat(e *Error) string {
	if e == nil {
		return "null"
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	_, _ = buf.WriteString("{")

	// ErrorType
	_, _ = buf.WriteString("\"error_type\":")
	_, _ = buf.WriteString("\"" + e.errorType.String() + "\",")

	// Severity
	_, _ = buf.WriteString("\"severity\":")
	_, _ = buf.WriteString("\"" + e.severity.String() + "\",")

	// Operations
	_, _ = buf.WriteString("\"operations\":[")
	ops := e.Operations()
	if len(ops) > 0 {
		op0 := ops[0].String()
		_, _ = buf.WriteString("\"" + op0 + "\"")
		for _, s := range ops[1:] {
			_, _ = buf.WriteString(",")
			opN := s.String()
			_, _ = buf.WriteString("\"" + opN + "\"")
		}
	}
	_, _ = buf.WriteString("],")

	// ContextInfo
	_, _ = buf.WriteString("\"context\":")
	data, _ := json.Marshal(e.ContextInfo())
	_, _ = buf.Write(data)
	_, _ = buf.WriteString(",")

	// Msg
	_, _ = buf.WriteString("\"msg\":")
	_, _ = buf.WriteString("\"" + e.TranslateMsg() + "\"")

	_, _ = buf.WriteString("}")

	return buf.String()
}

func StringFormat(e *Error) string {
	if e == nil {
		return ""
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	writeDelim := false

	if e.ErrorType() != "" {
		_, _ = buf.WriteString("[" + e.errorType.String() + "]")
		writeDelim = true
	}

	if e.Severity() > log.SeverityUnknown {
		_, _ = buf.WriteString("[" + e.severity.String() + "]")
		writeDelim = true
	}

	ops := e.Operations()
	if len(ops) > 0 {
		writeDelim = true
		_, _ = buf.WriteString("[")
		op0 := ops[0].String()
		_, _ = buf.WriteString(op0)
		for _, s := range ops[1:] {
			_, _ = buf.WriteString(",")
			opN := s.String()
			_, _ = buf.WriteString(opN)
		}
		_, _ = buf.WriteString("]")
	}

	ctxs := e.ContextInfo()
	if len(ctxs) > 0 {
		writeDelim = true
		_, _ = buf.WriteString("<")
		ctxskeys := make([]string, 0, len(ctxs))
		for i := range ctxs {
			ctxskeys = append(ctxskeys, i)
		}
		sort.Strings(ctxskeys)
		_, _ = buf.WriteString(fmt.Sprintf("%s:%v", ctxskeys[0], ctxs[ctxskeys[0]]))
		for _, i := range ctxskeys[1:] {
			_, _ = buf.WriteString(",")
			_, _ = buf.WriteString(fmt.Sprintf("%s:%v", i, ctxs[i]))
		}
		_, _ = buf.WriteString(">")
	}

	msg := e.Msg()
	if writeDelim && msg != "" {
		_, _ = buf.WriteString(" -- ")
	}

	if msg != "" {
		_, _ = buf.WriteString(e.TranslateMsg())
	}

	return buf.String()
}

//

// Format поддержка fmt.Printf()
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'c':
		fmt.Fprintf(s, "%v\n", e.ContextInfo())

	case 'o':
		fmt.Fprintf(s, "%v\n", e.Operations())

	case 'l':
		_, _ = io.WriteString(s, e.Severity().String())

	case 't':
		_, _ = io.WriteString(s, e.ErrorType().String())

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
