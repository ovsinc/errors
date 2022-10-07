package errors

// Options опции из параметра ошибки.
type Options func(e *Error)

// Msg

// SetMsg установит сообщение об ошибке, указанное в виде строки.
func SetMsg(msg string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = []byte(msg)
	}
}

// ID

// SetID установит ID ошибки.
func SetID(id string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = []byte(id)
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
		e.operation = []byte(o)
	}
}

// Error type

// SetErrorType установит тип, как строку.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetErrorType(et string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.errorType = []byte(et)
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
		e.contextInfo = append(e.contextInfo, struct{ Key, Value []byte }{[]byte(key), []byte(value)})
	}
}
