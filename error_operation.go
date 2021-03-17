package errors

// AppendOperations добавить операции.
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то указанные в аргументе будет добавлены к уже имеющимся.
func AppendOperations(ops ...string) Options {
	return func(e *Error) {
		if e == nil || ops == nil {
			return
		}
		e.operations = append(
			make([]string, 0, len(e.operations)+len(ops)),
			e.operations...,
		)
		e.operations = append(e.operations, ops...)
	}
}

// SetOperations установить операции
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetOperations(ops ...string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.operations = append(make([]string, 0, len(ops)), ops...)
	}
}

//

// Operations вернет список операций.
func (e *Error) Operations() []string {
	return e.operations
}
