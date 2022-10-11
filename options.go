package errors

import (
	"reflect"
	"unsafe"
)

// Options опции из параметра ошибки.
type Options func(e *Error)

// Msg

// SetMsg установит сообщение об ошибке, указанное в виде строки.
func SetMsg(msg string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = s2b(msg)
	}
}

// ID

// SetID установит ID ошибки.
func SetID(id string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = s2b(id)
	}
}

// Operation

// SetOperation установит операцию, как строку.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetOperation(o string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.operation = s2b(o)
	}
}

// Error type

// SetErrorType установит тип.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetErrorType(et errType) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.errorType = et
	}
}

// Context Info

// SetContextInfo установить CtxMap.
func SetContextInfo(ctxinf CtxKV) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.contextInfo = ctxinf
	}
}

// AppendContextInfo добавить в имеющийся CtxMap значение value по ключу key.
// Если CtxMap в *Error не установлен, то он будет предварительно установлен.
func AppendContextInfo(key string, value string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		if e.contextInfo == nil {
			e.contextInfo = make(CtxKV, 0, 6)
		}
		e.contextInfo = append(e.contextInfo, struct{ Key, Value []byte }{s2b(key), s2b(value)})
	}
}

//
// from https://github.com/valyala/fastjson/blob/master/util.go
//

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
	strh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Data = strh.Data
	sh.Len = strh.Len
	sh.Cap = strh.Len
	return b
}
