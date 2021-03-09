package errors

import (
	origerrors "errors"

	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type TranslateMessage struct {
	ID           string
	Localizer    *i18n.Localizer
	TemplateData map[string]interface{}
	PluralCount  interface{}
}

type translateMap map[string]*TranslateMessage

// SetLang установить используемый язык
func SetLang(lang string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.lang = lang
	}
}

// AddTranslatedMessage установить сообщение переревода и локализатор
func AddTranslatedMessage(lang string, message *TranslateMessage) Options {
	return func(e *Error) {
		if e == nil {
			return
		}

		if e.translateMap == nil {
			e.translateMap = make(translateMap)
		}
		e.translateMap[lang] = message
	}
}

//

func (e *Error) Lang() string {
	return e.lang
}

func (e *Error) TranslatedMessage(lang ...string) *TranslateMessage {
	if len(lang) == 0 {
		return e.translateMap[e.lang]
	}
	return e.translateMap[lang[0]]
}

//

var (
	ErrNoLocalizer        = origerrors.New("no localizer config for this lang")
	ErrNoTranslateContext = origerrors.New("no translate context for this lang")
)

func (e *Error) trans(s string) (string, error) {
	msgCtx := e.translateMap[e.lang]
	if msgCtx == nil {
		return s, ErrNoTranslateContext
	}

	localizer := msgCtx.Localizer
	if localizer == nil {
		return s, ErrNoLocalizer
	}

	msg, _, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID:    msgCtx.ID,
		TemplateData: msgCtx.TemplateData,
		PluralCount:  msgCtx.PluralCount,
	})
	if err != nil {
		return s, err
	}

	return msg, nil
}

func (e *Error) TranslateMsg() string {
	return e.Translate(e.msg)
}

func (e *Error) Translate(s string) string {
	s, _ = e.trans(s)
	return s
}

// helper

func Translate(err error, lang string) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(*Error); ok { // nolint:errorlint
		err = e.WithOptions(SetLang(lang))
	}

	return err.Error()
}
