package errors

import i18n "github.com/nicksnyder/go-i18n/v2/i18n"

// Options опции из параметра ошибки.
type Options func(e *Error)

// SetFormatFn установит пользовательскую функцию-форматирования.
func SetFormatFn(fn FormatFn) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.formatFn = fn
	}
}

// Msg

// SetMsgBytes установит сообщение об ошибке, указаннов в виде []byte.
func SetMsgBytes(msg []byte) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = NewObjectFromBytes(msg)
	}
}

// SetMsg установит сообщение об ошибке, указанное в виде строки.
func SetMsg(msg string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = NewObjectFromString(msg)
	}
}

// ID

// SetID установит ID ошибки.
func SetID(id string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = NewObjectFromString(id)
	}
}

// SetIDBytes установит ID ошибки.
func SetIDBytes(id []byte) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = NewObjectFromBytes(id)
	}
}

// Operation

// SetOperations установить операции, указанные как строки.
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetOperation(o string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.operation = NewObjectFromString(o)
	}
}

// SetOperationsBytes установить операции, указанные как []byte.
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetOperationsBytes(o []byte) Options {
	return func(e *Error) {
		if e == nil || len(o) == 0 {
			return
		}
		e.operation = NewObjectFromBytes(o)
	}
}

// Translate

// SetTranslateContext установит контекст переревода
func SetTranslateContext(tctx *TranslateContext) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.translateContext = tctx
	}
}

// SetLocalizer установит локализатор.
// Этот локализатор будет использован для данной ошибки даже,
// если был установлен DefaultLocalizer.
func SetLocalizer(localizer *i18n.Localizer) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.localizer = localizer
	}
}
