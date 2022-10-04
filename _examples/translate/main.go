package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/ovsinc/errors"
	"golang.org/x/text/language"
)

const (
	myName        = "John Snow"
	messagesCount = 5
)

func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("./testdata/active.ru.toml")
	localizer := i18n.NewLocalizer(bundle, "ru")

	err := errors.NewWith(
		errors.SetMsg("fallback message"),
		errors.SetID("ErrEmailsUnreadMsg"),
	)

	msg, _ := errors.Translate(err,
		localizer,
		&errors.TranslateContext{
			TemplateData: map[string]interface{}{
				"Name":        myName,
				"PluralCount": messagesCount,
			},
			PluralCount: messagesCount,
		},
	)
	fmt.Println(msg)

	errors.DefaultLocalizer = localizer

	eunknown := errors.NewWith(
		errors.SetMsg("fallback unknown message"),
		errors.SetID("ErrUnknownErrorMsg"),
	)

	fmt.Printf("%+s\n", eunknown)
}
