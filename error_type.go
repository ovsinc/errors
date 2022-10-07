package errors

import (
	"net/http"
)

//go:generate stringer -type=errType
type errType int

const (
	_ errType = iota
	// Unknown is error type for unknown system error. Default error type.
	Unknown
	// Internal is error type for when there is an internal system error. e.g. Database errors
	Internal
	// Validation is error type for when there is a validation error. e.g. invalid email address
	Validation
	// InputBody is error type for when an input data type error. e.g. invalid JSON
	InputBody
	// Duplicate is error type for when there's duplicate content
	Duplicate
	// Unauthenticated is error type when trying to access an authenticated API without authentication
	Unauthenticated
	// Unauthorized is error type for when there's an unauthorized access attempt
	Unauthorized
	// Empty is error type for when an expected non-empty resource, is empty
	Empty
	// NotFound is error type for an expected resource is not found e.g. user ID not found
	NotFound
	// MaximumAttempts is error type for attempting the same action more than allowed
	MaximumAttempts
	// SubscriptionExpired is error type for when a user's 'paid' account has expired
	SubscriptionExpired
	// DownstreamDependencyTimedout is error type for when a request to a downstream dependent service times out
	DownstreamDependencyTimedout
)

var defaultErrType = Internal

func ParseErrType(s string) errType {
	t := Unknown

	switch s {
	case "Validation":
		t = Validation

	case "InputBody":
		t = InputBody

	case "Duplicate":
		t = Duplicate

	case "Unauthenticated":
		t = Unauthenticated

	case "Unauthorized":
		t = Unauthorized

	case "Empty":
		t = Empty

	case "Internal":
		t = Internal

	case "NotFound":
		t = NotFound

	case "MaximumAttempts":
		t = MaximumAttempts

	case "SubscriptionExpired":
		t = SubscriptionExpired
	}

	return t
}

// HTTPStatusCode is a convenience method used to get the appropriate HTTP response status code for the respective error type
func (et errType) HTTPStatusCode() int {
	status := http.StatusTeapot

	switch et {
	case Validation:
		status = http.StatusUnprocessableEntity
	case InputBody:
		status = http.StatusBadRequest
	case Duplicate:
		status = http.StatusConflict
	case Unauthenticated:
		status = http.StatusUnauthorized
	case Unauthorized:
		status = http.StatusForbidden
	case Empty:
		status = http.StatusGone
	case NotFound:
		status = http.StatusNotFound
	case Internal:
		status = http.StatusInternalServerError
	case MaximumAttempts:
		status = http.StatusTooManyRequests
	case SubscriptionExpired:
		status = http.StatusPaymentRequired
	case Unknown:
		status = http.StatusTeapot
	}

	return status
}

func (i errType) Bytes() []byte {
	if i < 0 || i >= errType(len(_errType_index)-1) {
		return []byte("Unknown")
	}
	return []byte(_errType_name[_errType_index[i]:_errType_index[i+1]])
}

//

func getErrType(err error) (errType, bool) {
	errType := defaultErrType
	ok := false

	if e, eok := err.(*Error); eok {
		errType = e.ErrorType()
		ok = true
	}

	return errType, ok
}

func HTTPStatusCode(err error) (int, bool) {
	errType, ok := getErrType(err)
	return errType.HTTPStatusCode(), ok
}

func HTTPStatusCodeMessage(err error) (int, string, bool) {
	errType, ok := getErrType(err)
	return errType.HTTPStatusCode(), errType.String(), ok
}

func StatusMessage(err error) (string, bool) {
	errType, ok := getErrType(err)
	return errType.String(), ok
}

//

func errWithType(eType errType, ops ...Options) *Error {
	op := make([]Options, 0, len(ops)+1)
	op = append(op, ops...)
	op = append(op, SetErrorType(eType))
	return NewWith(op...)
}

func IternalErrWith(ops ...Options) *Error {
	return errWithType(Internal, ops...)
}

func IternalErr(s string) *Error {
	return IternalErrWith(SetMsg(s))
}

func ValidationErrWith(ops ...Options) *Error {
	return errWithType(Validation, ops...)
}

func ValidationErr(s string) *Error {
	return ValidationErrWith(SetMsg(s))
}

func InputBodyErrWith(ops ...Options) *Error {
	return errWithType(InputBody, ops...)
}

func InputBodyErr(s string) *Error {
	return InputBodyErrWith(SetMsg(s))
}

func UnauthenticatedErrWith(ops ...Options) *Error {
	return errWithType(Unauthenticated, ops...)
}

func UnauthenticatedErr(s string) *Error {
	return UnauthenticatedErrWith(SetMsg(s))
}

func UnauthorizedErrWith(ops ...Options) *Error {
	return errWithType(Unauthorized, ops...)
}

func UnauthorizedErr(s string) *Error {
	return UnauthenticatedErrWith(SetMsg(s))
}

func DuplicateErrWith(ops ...Options) *Error {
	return errWithType(Duplicate, ops...)
}

func DuplicateErr(s string) *Error {
	return DuplicateErrWith(SetMsg(s))
}

func EmptyErrWith(ops ...Options) *Error {
	return errWithType(Empty, ops...)
}

func EmptyErr(s string) *Error {
	return EmptyErrWith(SetMsg(s))
}

func NotFoundErrWith(ops ...Options) *Error {
	return errWithType(NotFound, ops...)
}

func NotFoundErr(s string) *Error {
	return NotFoundErrWith(SetMsg(s))
}

func MaximumAttemptsErrWith(ops ...Options) *Error {
	return errWithType(MaximumAttempts, ops...)
}

func MaximumAttemptsErr(s string) *Error {
	return MaximumAttemptsErrWith(SetMsg(s))
}

func SubscriptionExpiredErrWith(ops ...Options) *Error {
	return errWithType(SubscriptionExpired, ops...)
}

func SubscriptionExpiredErr(s string) *Error {
	return SubscriptionExpiredErrWith(SetMsg(s))
}

func DownstreamDependencyTimedoutErrWith(ops ...Options) *Error {
	return errWithType(DownstreamDependencyTimedout, ops...)
}

func DownstreamDependencyTimedoutErr(s string) *Error {
	return DownstreamDependencyTimedoutErrWith(SetMsg(s))
}
