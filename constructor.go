package errors

// New конструктор на необязательных параметрах
// * ops ...Options -- параметризация через функции-парметры.
// См. options.go
//
// ** *Error
func New(i interface{}) *Error {
	var msg string
	switch t := i.(type) {
	case string:
		msg = t
	case error:
		msg = t.Error()
	case interface{ String() string }:
		msg = t.String()
	case func() string:
		msg = t()
	}
	return NewWith(SetMsg(msg))
}

func NewWith(ops ...Options) *Error {
	e := Error{}
	for _, op := range ops {
		op(&e)
	}
	return &e
}

// NewLog конструктор *Error, как и New,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func NewWithLog(i interface{}) *Error {
	e := New(i)
	e.Log()
	return e
}

//

// Combine создаст цепочку ошибок из ошибок ...errors.
// Допускается использование `nil` в аргументах.
func Combine(errors ...error) Multierror {
	return fromSlice(errors)
}

// Wrap обернет ошибку `left` ошибкой `right`, получив цепочку.
// Допускается использование `nil` в одном из аргументов.
func Wrap(left error, right error) Multierror {
	return fromSlice([]error{left, right})
}

// append err to []*Error
// errors must not be nil
func appendError(errors []*Error, err interface{}) []*Error {
	switch t := err.(type) {
	case nil:
		return nil

	case *Error:
		return append(errors, t)

	case Multierror:
		return append(errors, t.Errors()...)

	case error:
		return append(errors, New(t.Error()))
	}

	return errors
}

// fromSlice converts the given list of errors into a single error.
func fromSlice(errors []error) Multierror {
	nonNilErrs := make([]*Error, 0)
	for _, err := range errors {
		if err == nil {
			continue
		}
		nonNilErrs = appendError(nonNilErrs, err)
	}

	last := 0
	len := len(nonNilErrs)
	if len > 0 {
		last = len - 1
	}

	return &multiError{
		errors: nonNilErrs,
		len:    len,
		last:   last,
	}
}

//

// Options опции из параметра ошибки.
type Options func(e *Error)

// Msg

// SetMsg установит сообщение об ошибке, указанное в виде строки.
func SetMsg(msg string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = NewMsgFromString(msg)
	}
}

// ID

// SetID установит ID ошибки.
func SetID(id string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = NewIDFromString(id)
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
		e.operation = NewOperationFromString(o)
	}
}

func WithCaller(depth ...CallDepth) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		d := DefaultCallDepth
		if len(depth) > 0 {
			d = depth[0]
		}
		e.caller = Caller(d)
	}
}

// Error type

// SetErrorType установит тип, как строку.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetErrorType(o string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.operation = NewErrorTypeFromString(o)
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
func SetLocalizer(localizer Localizer) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.localizer = localizer
	}
}

// Context Info

// SetContextInfo установить CtxMap.
func SetContextInfo(ctxinf CtxMap) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.contextInfo = ctxinf
	}
}

// AppendContextInfo добавить в имеющийся CtxMap значение value по ключу key.
// Если CtxMap в *Error не установлен, то он будет предварительно установлен.
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
