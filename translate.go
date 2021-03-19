package errors

import (
	origerrors "errors"
	"io"

	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// DefaultLocalizer локализатор по-умолчанию.
// Для каждой ошибки можно переопределить локализатор.
var DefaultLocalizer *i18n.Localizer //nolint:gochecknoglobals

// TranslateContext контекст перевода. Не является обязательным для корректного перевода.
type TranslateContext struct {
	// TemplateData - map для замены в шаблоне
	TemplateData map[string]interface{}
	// PluralCount признак множественности.
	// Может иметь значение nil или число.
	PluralCount interface{}
	// DefaultMessage сообщение, которое будет использовано при ошибке перевода.
	DefaultMessage *i18n.Message
}

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

//

// TranslateContext вернет *TranslateContext.
func (e *Error) TranslateContext() *TranslateContext {
	return e.translateContext
}

// Localizer вернет локализатор *i18n.Localizer.
func (e *Error) Localizer() *i18n.Localizer {
	return e.localizer
}

//

var errNoLocalizer = origerrors.New("no localizer config for this lang")

// WriteTranslateMsg запишет перевод сообщения ошибки в буфер.
func (e *Error) WriteTranslateMsg(w io.Writer) {
	_ = e.writeTranslateMsg(w)
}

func (e *Error) writeTranslateMsg(w io.Writer) error {
	s := e.msg

	if len(s) == 0 {
		return nil
	}

	var localizer *i18n.Localizer
	switch {
	case e.localizer != nil:
		localizer = e.localizer
	case DefaultLocalizer != nil:
		localizer = DefaultLocalizer
	}

	if localizer == nil {
		_, _ = io.WriteString(w, s)
		return errNoLocalizer
	}

	i18nConf := i18n.LocalizeConfig{
		MessageID: e.id,
	}
	if e.translateContext != nil {
		i18nConf.DefaultMessage = e.translateContext.DefaultMessage
		i18nConf.PluralCount = e.translateContext.PluralCount
		i18nConf.TemplateData = e.translateContext.TemplateData
	}

	msg, _, err := e.localizer.LocalizeWithTag(&i18nConf)
	if err != nil {
		_, _ = io.WriteString(w, s)
		return err
	}

	_, _ = io.WriteString(w, msg)

	return nil
}
