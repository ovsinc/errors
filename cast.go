package errors

type Caster interface {
	Cast(interface{})
}

type ErrorOut struct {
	Message   string                 `json:"message"`
	Operation string                 `json:"operation"`
	ID        string                 `json:"id"`
	Context   map[string]interface{} `json:"context"`
}

func (out *ErrorOut) Cast(e interface{}) {
	if x, ok := e.(*Error); ok { //nolint:errorlint
		out.ID = x.ID().String()
		out.Context = x.ContextInfo()
		out.Message = x.TranslateMsg()
		out.Operation = x.Operation().String()
	}
}

//

type MultierrorOut struct {
	Count  int         `json:"count"`
	Errors []*ErrorOut `json:"errors"`
	Title  string      `json:"title"`
}

func (meo *MultierrorOut) Cast(e interface{}) {
	if x, ok := e.(Multierror); ok { //nolint:errorlint
		meo.Count = x.Len()
		meo.Errors = make([]*ErrorOut, 0, meo.Count)
		for _, e := range x.Errors() {
			var out ErrorOut
			e.CastTo(&out)
			meo.Errors = append(meo.Errors, &out)
		}
	}
}

//

func (e *Error) CastTo(c Caster) {
	if e != nil && c != nil {
		c.Cast(e)
	}
}

func (merr *multiError) CastTo(c Caster) {
	if merr != nil && c != nil {
		c.Cast(merr)
	}
}

//

func simpleCast(err error) (*Error, bool) {
	e, ok := err.(*Error) //nolint:errorlint
	return e, ok
}

// Cast преобразует тип error в *Error
// Если error не соответствует *Error, то будет создан *Error с сообщением err.Error().
// Для err == nil, вернется nil.
func Cast(err error) *Error {
	if e, ok := simpleCast(err); ok {
		return e
	}
	return New(err.Error())
}
