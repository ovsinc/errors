package errors

import (
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// DefaultLocalizer локализатор по-умолчанию.
// Для каждой ошибки можно переопределить локализатор.
var DefaultLocalizer Localizer //nolint:gochecknoglobals

var (
	ErrNoLocalizer = New("localizer is no set")
	ErrNotError    = New("not *Error")
)

type Localizer interface {
	Localize(*i18n.LocalizeConfig) (string, error)
}

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

func DefaultTranslate(e error) string {
	msg, _ := Translate(e, nil, nil)
	return msg
}

// Translate вернет перевод сообщения ошибки.
// Если не удастся выполнить перевод, вернет оригинальное сообщение.
func Translate(e error, l Localizer, tctx *TranslateContext) (string, error) {
	err, ok := e.(*Error) //nolint:errorlint
	var loc Localizer

	switch {
	case !ok:
		return e.Error(), ErrNotError

	case l != nil:
		loc = l

	case DefaultLocalizer != nil:
		loc = DefaultLocalizer

	default:
		// no localizer
		return err.Msg(), ErrNoLocalizer
	}

	i18nConf := i18n.LocalizeConfig{
		MessageID: err.ID(),
	}
	if tctx != nil {
		i18nConf.DefaultMessage = tctx.DefaultMessage
		i18nConf.PluralCount = tctx.PluralCount
		i18nConf.TemplateData = tctx.TemplateData
	}

	if msg, err := loc.Localize(&i18nConf); err == nil {
		return msg, nil
	}

	return err.Msg(), nil
}
