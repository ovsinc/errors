package errors

// NewOperation констуруктор Operation на основе значения s.
func NewOperation(s string) Operation {
	return Operation(s)
}

// Operation тип операции.
type Operation string

// String вернет строковое представление операции.
func (o Operation) String() string {
	return string(o)
}

//

// AppendOperations добавить операции.
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то указанные в аргументе будет добавлены к уже имеющимся.
func AppendOperations(ops ...Operation) Options {
	return func(e *Error) {
		if e == nil || ops == nil {
			return
		}
		e.operations = append(
			make([]Operation, 0, len(e.operations)+len(ops)),
			e.operations...,
		)
		e.operations = append(e.operations, ops...)
	}
}

// SetOperations установить операции
// Можно указать произвольное количество.
// Если в *Error уже были записаны операции,
// то они будут заменены на указанные в аргументе ops.
func SetOperations(ops ...Operation) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.operations = append(make([]Operation, 0, len(ops)), ops...)
	}
}

//

// Operations вернет список операций.
func (e *Error) Operations() []Operation {
	return e.operations
}
