package errors

import (
	"runtime"
	"strconv"
	"strings"
)

//
// From https://github.com/go-kit/kit/blob/master/log/value.go
//

type CallDepth int

const (
	DefaultCallDepth CallDepth = 3
	RuntimeCallDepth CallDepth = 2
	HandlerCallDepth CallDepth = 1
)

var (
	// DefaultCaller is a Valuer that returns then CallInfo where the Log
	// method was invoked. It can only be used with log.With.
	DefaultCaller = Caller(DefaultCallDepth) //nolint:gochecknoglobals

	RuntimeCaller = Caller(RuntimeCallDepth) //nolint:gochecknoglobals

	// HandlerCaller is a Valuer that returns then CallInfo where the
	// method was invoked. It can only be used within handler method.
	HandlerCaller = Caller(HandlerCallDepth) //nolint:gochecknoglobals
)

// Caller returns a Valuer that returns a CallInfo from a specified depth
// in the callstack. Users will probably want to use DefaultCaller.
func Caller(depth CallDepth) func() Objecter {
	return func() Objecter {
		pc, file, line, _ := runtime.Caller(int(depth))

		idx := strings.LastIndexByte(file, '/')

		// using idx+1 below handles both of following cases:
		// idx == -1 because no "/" was found, or
		// idx >= 0 and we want to start at the character after the found "/".
		fpos := file[idx+1:] + ":" + strconv.Itoa(line)

		fn := runtime.FuncForPC(pc)
		if fn != nil {
			fname := fn.Name()
			nameidx := strings.LastIndexByte(fname, '.')
			funcName := fname[nameidx+1:]
			fpos = fpos + ": " + funcName + "()"

		}
		return NewFileLineFromString(fpos)
	}
}
