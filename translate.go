package errors

import (
	origerrors "errors"

	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

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

func (e *Error) trans(s string) (string, error) {
	if e.localizer == nil {
		return s, errNoLocalizer
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
		return s, err
	}

	return msg, nil
}

// translateMsg выполнит перевод сообщения об ошибке.
// Метод в случае неудачи перевода вернет оригинальное сообщение.
func (e *Error) translateMsg() string {
	s, _ := e.trans(e.msg)
	return s
}
