package errors

import (
	"runtime"
	"strconv"
	"strings"
)

type CallInfo struct {
	FilePosition string
	FuncName     string
	FullFuncName string
}

//
// From https://github.com/go-kit/kit/blob/master/log/value.go
//

// Caller returns a Valuer that returns a CallInfo from a specified depth
// in the callstack. Users will probably want to use DefaultCaller.
func Caller(depth int) func() CallInfo {
	return func() (ret CallInfo) {
		pc, file, line, _ := runtime.Caller(depth)

		idx := strings.LastIndexByte(file, '/')
		// using idx+1 below handles both of following cases:
		// idx == -1 because no "/" was found, or
		// idx >= 0 and we want to start at the character after the found "/".

		ret.FilePosition = file[idx+1:] + ":" + strconv.Itoa(line)

		fn := runtime.FuncForPC(pc)
		if fn != nil {
			ret.FullFuncName = fn.Name()
			nameidx := strings.LastIndexByte(ret.FullFuncName, '.')
			ret.FuncName = ret.FullFuncName[nameidx+1:]

		}

		return
	}
}

// DefaultCaller is a Valuer that returns then CallInfo where the Log
// method was invoked. It can only be used with log.With.
var DefaultCaller = Caller(3) //nolint:gochecknoglobals

// HandlerCaller is a Valuer that returns then CallInfo where the
// method was invoked. It can only be used within handler method.
var HandlerCaller = Caller(1) //nolint:gochecknoglobals
