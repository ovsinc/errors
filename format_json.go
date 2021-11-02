package errors

import (
	"io"
	"strconv"

	json "github.com/goccy/go-json"
)

// JSONFormat функция форматирования вывода сообщения *Error в JSON.
func JSONFormat(buf io.Writer, e *Error) {
	if e == nil {
		_, _ = io.WriteString(buf, "null")
		return
	}

	_, _ = io.WriteString(buf, "{")

	// ID
	_, _ = io.WriteString(buf, "\"id\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = buf.Write(e.ID().Bytes())
	_, _ = io.WriteString(buf, "\",")

	// Operation
	_, _ = io.WriteString(buf, "\"operation\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = buf.Write(e.Operation().Bytes())
	_, _ = io.WriteString(buf, "\",")

	// ContextInfo
	_, _ = io.WriteString(buf, "\"context\":")
	cxtInfo := e.ContextInfo()
	if len(cxtInfo) > 0 {
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "")
		_ = enc.Encode(e.ContextInfo())
	} else {
		_, _ = io.WriteString(buf, "null")
	}
	_, _ = io.WriteString(buf, ",")

	// Msg
	_, _ = io.WriteString(buf, "\"msg\":")
	_, _ = io.WriteString(buf, "\"")
	if len(e.Msg().Bytes()) > 0 {
		_, _ = e.WriteTranslateMsg(buf)
	}
	_, _ = io.WriteString(buf, "\"")

	_, _ = io.WriteString(buf, "}")
}

// JSONMultierrFuncFormat функция форматирования вывода сообщения для multierr в виде JSON.
func JSONMultierrFuncFormat(w io.Writer, es []*Error) {
	if len(es) == 0 {
		_, _ = io.WriteString(w, "null")
	}

	_, _ = io.WriteString(w, "{")

	_, _ = io.WriteString(w, "\"count\":")
	_, _ = io.WriteString(w, strconv.Itoa(len(es)))
	_, _ = io.WriteString(w, ",")

	_, _ = io.WriteString(w, "\"messages\":")
	_, _ = io.WriteString(w, "[")
	switch len(es) {
	case 1:
		JSONFormat(w, es[0])
	default:
		JSONFormat(w, es[0])
		for _, e := range es[1:] {
			_, _ = io.WriteString(w, ",")
			JSONFormat(w, e)
		}
	}
	_, _ = io.WriteString(w, "]")

	_, _ = io.WriteString(w, "}")
}
