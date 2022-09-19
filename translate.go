package errors

import (
	"io"

	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/valyala/bytebufferpool"
)

// DefaultLocalizer локализатор по-умолчанию.
// Для каждой ошибки можно переопределить локализатор.
var DefaultLocalizer Localizer //nolint:gochecknoglobals

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

// TranslateMsg вернет перевод сообщения ошибки.
// Если не удастся выполнить перевод, вернет оригинальное сообщение.
func (e *Error) TranslateMsg() string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	_, _ = e.WriteTranslateMsg(buf)

	return buf.String()
}

// WriteTranslateMsg запишет перевод сообщения ошибки в буфер.
// Если не удастся выполнить перевод в буфер w будет записано оригинальное сообщение.
func (e *Error) WriteTranslateMsg(w io.Writer) (int, error) {
	var loc Localizer
	switch {
	case e.msg == nil:
		return 0, nil
	case e.localizer != nil:
		loc = e.localizer
	case DefaultLocalizer != nil:
		loc = DefaultLocalizer
	default:
		// no localizer
		return e.Msg().Write(w)
	}

	i18nConf := &i18n.LocalizeConfig{
		MessageID: e.id.String(),
	}
	if e.translateContext != nil {
		i18nConf.DefaultMessage = e.translateContext.DefaultMessage
		i18nConf.PluralCount = e.translateContext.PluralCount
		i18nConf.TemplateData = e.translateContext.TemplateData
	}

	str, err := loc.Localize(i18nConf)
	// fallback
	if err != nil {
		return e.Msg().Write(w)
	}

	return io.WriteString(w, str)
}
