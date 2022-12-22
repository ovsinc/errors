package errors

import (
	"reflect"
	"unsafe"
)

// Options опции-параметры ошибки.
type Options func(e *Error)

// Msg

// SetMsg строка. Установит сообщение об ошибке.
func SetMsg(msg string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = msg
	}
}

// ID

// SetID, строка. Установит ID ошибки.
func SetID(id string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = id
	}
}

// Operation

// SetOperation, строка. Установит имя операции.
func SetOperation(o string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.operation = o
	}
}

// Error type

// SetErrorType, errType (enum). Установит тип.
func SetErrorType(et errType) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.errorType = et
	}
}

// Context Info

// SetContextInfo, CtxKV. Установит контекст.
func SetContextInfo(ctxinf CtxKV) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.contextInfo = ctxinf
	}
}

// AppendContextInfo, key, val - строки. Добавит в имеющийся CtxKV значение value по ключу key.
// CtxKV будет инициализирован, если ранее этого не было сделано.
func AppendContextInfo(key string, value interface{}) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		if e.contextInfo == nil {
			e.contextInfo = make(CtxKV, 0, 6)
		}
		e.contextInfo = append(e.contextInfo, struct {
			Key   string
			Value interface{}
		}{key, value})
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
