package main

import (
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"

	"github.com/ovsinc/errors"
)

var (
	ErrBadContent = errors.NewWith(
		errors.SetID(EBadContentMsg.ID),
		errors.SetErrorType(EInputBody.String()),
		errors.SetMsg(EBadContentMsg.Other),
	)
)

type ErrorMessage struct {
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Type      string `json:"type"`
	ID        string `json:"id"`
}

type Response struct {
	Data   interface{}   `json:"data"`
	HasErr bool          `json:"has_error"`
	Error  *ErrorMessage `json:"error"`
}

type Handler struct {
	uc        *UserUC
	localizer *translater
}

func NewHandler(uc *UserUC, localizer *translater) *Handler {
	return &Handler{
		uc:        uc,
		localizer: localizer,
	}
}

func (h *Handler) Register(app *echo.Echo) {
	app.GET("/user/:id", h.GetUser)
	app.DELETE("/user/:id", h.DeleteUser)
	app.POST("/user", h.NewUser)
}

var ErrErrHandle = errors.New("can`t error handling")

func (h *Handler) prepareMsg(err *errors.Error, lang string) (int, *Response) {
	return ParseIntErr(string(err.ErrorType())).Status(),
		&Response{
			Data:   nil,
			HasErr: true,
			Error: &ErrorMessage{
				Message:   h.localizer.TranslateError(lang, err),
				Operation: string(err.Operation()),
				Type:      string(err.ErrorType()),
				ID:        string(err.ID()),
			},
		}
}

func (h *Handler) errHandle(c echo.Context, err error) error {
	if err == nil {
		return c.JSON(
			http.StatusOK,
			Response{
				HasErr: false,
			},
		)
	}

	errors.Log(err)

	var e *errors.Error
	switch {
	case errors.ContainsByID(err, EBadContentMsg.ID):
		e = errors.UnwrapByID(err, EBadContentMsg.ID)

	case errors.ContainsByID(err, EValidationMsg.ID):
		e = errors.UnwrapByID(err, EValidationMsg.ID)

	case errors.ContainsByID(err, ENotFoundMsg.ID):
		e = errors.UnwrapByID(err, ENotFoundMsg.ID)

	case errors.ContainsByID(err, EDuplicateMsg.ID):
		e = errors.UnwrapByID(err, EDuplicateMsg.ID)

	case errors.ContainsByID(err, EEmptyMsg.ID):
		e = errors.UnwrapByID(err, EEmptyMsg.ID)

	case errors.ContainsByID(err, EInternalMsg.ID):
		e = errors.UnwrapByID(err, EInternalMsg.ID)

	default:
		e = ErrDBInternal
	}

	code, msg := h.prepareMsg(e, c.Request().Header.Get("Accept-Language"))
	return c.JSON(code, msg)

}

func (h *Handler) GetUser(c echo.Context) error {
	usr, err := h.uc.GetUser(c.Request().Context(), c.Param("id"))
	if err == nil {
		return c.JSON(
			http.StatusOK,
			Response{
				Data:   usr,
				HasErr: false,
			})
	}

	return h.errHandle(c, err)
}

func (h *Handler) NewUser(c echo.Context) error {
	u := new(User)

	ctype := c.Request().Header.Get(echo.HeaderContentType)
	if !strings.HasPrefix(ctype, echo.MIMEApplicationJSON) {
		return h.errHandle(c,
			ErrBadContent.WithOptions(errors.SetOperation("handler.NewUser")))
	}

	if err := c.Bind(u); err != nil {
		return h.errHandle(c,
			errors.Wrap(ErrBadContent.WithOptions(errors.SetOperation("handler.NewUser")), err))
	}

	return h.errHandle(c,
		h.uc.NewUser(c.Request().Context(), u))
}

func (h *Handler) DeleteUser(c echo.Context) error {
	return h.errHandle(c,
		h.uc.DeleteUser(c.Request().Context(), c.Param("id")))
}
