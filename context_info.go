package errors

type CtxMap map[string]interface{}

// SetContextInfo ...
func SetContextInfo(ctxinf CtxMap) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.contextInfo = ctxinf
	}
}

// AppendContextInfo ...
func AppendContextInfo(key string, value interface{}) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		if e.contextInfo == nil {
			e.contextInfo = make(CtxMap)
		}
		e.contextInfo[key] = value
	}
}

//

// ContextInfo получить контекст ошибки
func (e *Error) ContextInfo() CtxMap {
	return e.contextInfo
}
