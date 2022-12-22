package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

//go:generate stringer -type=errType
type errType int

const (
	_ errType = iota

	// Unknown is error type for unknown system error. Default error type.

	// Неизвестный тип ошибки. Дефолтное значение.
	Unknown

	// Internal is error type for when there is an internal system error. e.g. Database errors

	// Internal внутренняя системная ошибка. Например, отказ базы данных.
	Internal

	// Validation is error type for when there is a validation error. e.g. invalid email address

	// Validation ошибка валидации. Например, не корректный email-адрес.
	Validation

	// InputBody is error type for when an input data type error. e.g. invalid JSON

	// InputBody ошибка обработки входных данных. Например, ошибка сериализации JSON.
	InputBody

	// Duplicate is error type for when there's duplicate content

	// Duplicate дубликат данных, нарушения уникальности.
	Duplicate

	// Unauthenticated is error type when trying to access an authenticated API without authentication

	// Unauthenticated для выполнения запроса требуется аутентфиикация.
	Unauthenticated

	// Unauthorized is error type for when there's an unauthorized access attempt

	// Unauthorized доступ запрещен, запрос не авторизован.
	Unauthorized

	// Empty is error type for when an expected non-empty resource, is empty

	// Empty запрос или не ответ не должен быть пустым.
	Empty

	// NotFound is error type for an expected resource is not found e.g. user ID not found

	// NotFound запрашиваемые данные не найдены. Например, пользователь с заданным ID не найден.
	NotFound

	// MaximumAttempts is error type for attempting the same action more than allowed

	// MaximumAttempts превышение числе разрешенных попуток выполнения одного и того же действия.
	MaximumAttempts

	// SubscriptionExpired is error type for when a user's 'paid' account has expired

	// SubscriptionExpired срок действия "оплаченой" подписки истек.
	SubscriptionExpired

	// DownstreamDependencyTimedout is error type for when a request to a downstream dependent service times out

	// DownstreamDependencyTimedout время ожидания выполнения запрос к нижестоящему сервису истек.
	DownstreamDependencyTimedout

	// Unavailable is error type for when server is unavailable.

	// Unavailable сервис не доступен.
	Unavailable
)

var defaultErrType = Unknown

// ParseErrType позволяет получить errType по названию.
func ParseErrType(s string) errType { //nolint:cyclop
	t := defaultErrType

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

	case "Unavailable":
		t = Unavailable
	}

	return t
}

// HTTPStatusCode is a convenience method used to get the appropriate HTTP response status code for the respective error type

// HTTPStatusCode позволяет конвертировать errType в HTTP status.
func (et errType) HTTPStatusCode() int { //nolint:cyclop
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
	case DownstreamDependencyTimedout:
		status = http.StatusRequestTimeout
	case Unavailable:
		status = http.StatusServiceUnavailable
	}

	return status
}

// GRPCStatusCode is a convenience method used to get the appropriate gRPC response code for the respective error type

// GRPCStatusCode позволяет конвертировать errType в gRPC status.
func (et errType) GRPCStatusCode() codes.Code { //nolint:cyclop
	status := codes.Unknown

	switch et {
	case NotFound:
		status = codes.NotFound
	case Duplicate:
		status = codes.AlreadyExists
	case Validation, InputBody, Empty:
		status = codes.InvalidArgument
	case Internal:
		status = codes.Internal
	case Unauthenticated:
		status = codes.Unauthenticated
	case Unauthorized:
		status = codes.PermissionDenied
	case MaximumAttempts, Unavailable, SubscriptionExpired:
		status = codes.Unavailable
	case DownstreamDependencyTimedout:
		status = codes.DeadlineExceeded
	}

	return status
}

//

// GetErrType получить errType из error.
// * in: error
// * out: t errType, ok bool
// Если error кастится на (*Error), то ok == true, и возвращается значение errType.
// В противном случае возвращается defaultErrType и false.
func GetErrType(err error) (errType, bool) {
	errType := defaultErrType
	ok := false

	if e, eok := err.(*Error); eok { //nolint:errorlint
		errType = e.ErrorType()
		ok = true
	}

	return errType, ok
}

// GRPCStatusCode получить gRPC статус из error.
// * in: error
// * out: t codes.Code, ok bool
// Если error кастится на (*Error), то ok == true, и возвращается значение codes.Code,
// соответсвующее errType.
// В противном случае возвращается codes.Unknown и false.
func GRPCStatusCode(err error) (codes.Code, bool) {
	errType, ok := GetErrType(err)
	return errType.GRPCStatusCode(), ok
}

// HTTPStatusCode получить HTTP статус из error.
// * in: error
// * out: t int, ok bool
// Если error кастится на (*Error), то ok == true, и возвращается значение,
// соответсвующее errType.
// В противном случае возвращается http.StatusTeapot и false.
func HTTPStatusCode(err error) (int, bool) {
	errType, ok := GetErrType(err)
	return errType.HTTPStatusCode(), ok
}

// GRPCStatusCodeMessage получить gRPC статус из error.
// * in: error
// * out: t codes.Code, s string, ok bool
// Если error кастится на (*Error), то ok == true, возвращается значение t,
// соответсвующее errType и текстовое предсталение errType.
// В противном случае возвращается codes.Unknown, "Unknown", false.
func GRPCStatusCodeMessage(err error) (codes.Code, string, bool) {
	errType, ok := GetErrType(err)
	return errType.GRPCStatusCode(), errType.String(), ok
}

// HTTPStatusCodeMessage получить HTTP статус из error.
// * in: error
// * out: t int, s string, ok bool
// Если error кастится на (*Error), то ok == true, возвращается значение t,
// соответсвующее errType и текстовое предсталение errType.
// В противном случае возвращается http.StatusTeapot, "Unknown", false.
func HTTPStatusCodeMessage(err error) (int, string, bool) {
	errType, ok := GetErrType(err)
	return errType.HTTPStatusCode(), errType.String(), ok
}

// StatusMessage получить строковое описание errType из error.
// * in: error
// * out: s string, ok bool
// Если error кастится на (*Error), то ok == true, возвращается значение t,
// соответсвующее errType и текстовое предсталение errType.
// В противном случае возвращается "Unknown", false.
func StatusMessage(err error) (string, bool) {
	errType, ok := GetErrType(err)
	return errType.String(), ok
}

//

func errWithType(eType errType, ops ...Options) *Error {
	op := make([]Options, 0, len(ops)+1)
	op = append(op, ops...)
	op = append(op, SetErrorType(eType))
	return NewWith(op...)
}

// Конструктор *Error c типом Internal.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func IternalErrWith(ops ...Options) *Error {
	return errWithType(Internal, ops...)
}

// Конструктор *Error c типом Internal.
// * s string -- сообщение ошибки.
// ** *Error
func IternalErr(s string) *Error {
	return IternalErrWith(SetMsg(s))
}

// Конструктор *Error c типом Validation.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func ValidationErrWith(ops ...Options) *Error {
	return errWithType(Validation, ops...)
}

// Конструктор *Error c типом Validation.
// * s string -- сообщение ошибки.
// ** *Error
func ValidationErr(s string) *Error {
	return ValidationErrWith(SetMsg(s))
}

// Конструктор *Error c типом InputBody.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func InputBodyErrWith(ops ...Options) *Error {
	return errWithType(InputBody, ops...)
}

// Конструктор *Error c типом InputBody.
// * s string -- сообщение ошибки.
// ** *Error
func InputBodyErr(s string) *Error {
	return InputBodyErrWith(SetMsg(s))
}

// Конструктор *Error c типом Duplicate.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func DuplicateErrWith(ops ...Options) *Error {
	return errWithType(Duplicate, ops...)
}

// Конструктор *Error c типом Duplicate.
// * s string -- сообщение ошибки.
// ** *Error
func DuplicateErr(s string) *Error {
	return DuplicateErrWith(SetMsg(s))
}

// Конструктор *Error c типом Unauthenticated.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func UnauthenticatedErrWith(ops ...Options) *Error {
	return errWithType(Unauthenticated, ops...)
}

// Конструктор *Error c типом Unauthenticated.
// * s string -- сообщение ошибки.
// ** *Error
func UnauthenticatedErr(s string) *Error {
	return UnauthenticatedErrWith(SetMsg(s))
}

// Конструктор *Error c типом Unauthorized.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func UnauthorizedErrWith(ops ...Options) *Error {
	return errWithType(Unauthorized, ops...)
}

// Конструктор *Error c типом Unauthorized.
// * s string -- сообщение ошибки.
// ** *Error
func UnauthorizedErr(s string) *Error {
	return UnauthenticatedErrWith(SetMsg(s))
}

// Конструктор *Error c типом Empty.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func EmptyErrWith(ops ...Options) *Error {
	return errWithType(Empty, ops...)
}

// Конструктор *Error c типом Empty.
// * s string -- сообщение ошибки.
// ** *Error
func EmptyErr(s string) *Error {
	return EmptyErrWith(SetMsg(s))
}

// Конструктор *Error c типом NotFound.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func NotFoundErrWith(ops ...Options) *Error {
	return errWithType(NotFound, ops...)
}

// Конструктор *Error c типом NotFound.
// * s string -- сообщение ошибки.
// ** *Error
func NotFoundErr(s string) *Error {
	return NotFoundErrWith(SetMsg(s))
}

// Конструктор *Error c типом MaximumAttempts.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func MaximumAttemptsErrWith(ops ...Options) *Error {
	return errWithType(MaximumAttempts, ops...)
}

// Конструктор *Error c типом MaximumAttempts.
// * s string -- сообщение ошибки.
// ** *Error
func MaximumAttemptsErr(s string) *Error {
	return MaximumAttemptsErrWith(SetMsg(s))
}

// Конструктор *Error c типом SubscriptionExpired.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func SubscriptionExpiredErrWith(ops ...Options) *Error {
	return errWithType(SubscriptionExpired, ops...)
}

// Конструктор *Error c типом SubscriptionExpired.
// * s string -- сообщение ошибки.
// ** *Error
func SubscriptionExpiredErr(s string) *Error {
	return SubscriptionExpiredErrWith(SetMsg(s))
}

// Конструктор *Error c типом DownstreamDependencyTimedout.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func DownstreamDependencyTimedoutErrWith(ops ...Options) *Error {
	return errWithType(DownstreamDependencyTimedout, ops...)
}

// Конструктор *Error c типом DownstreamDependencyTimedout.
// * s string -- сообщение ошибки.
// ** *Error
func DownstreamDependencyTimedoutErr(s string) *Error {
	return DownstreamDependencyTimedoutErrWith(SetMsg(s))
}

// Конструктор *Error c типом Unavailable.
// * ops ...Options -- параметризация через функции-парметры.
// ** *Error
func UnavailableErrWith(ops ...Options) *Error {
	return errWithType(Unavailable, ops...)
}

// Конструктор *Error c типом Unavailable.
// * s string -- сообщение ошибки.
// ** *Error
func UnavailableErr(s string) *Error {
	return UnavailableErrWith(SetMsg(s))
}
