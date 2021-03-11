package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/ovsinc/errors"
	"golang.org/x/text/language"
)

func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("./testdata/active.ru.toml")

	err := errors.New(
		"fallback message",
		errors.SetID("ErrEmailsUnreadMsg"),
		errors.SetErrorType(errors.NewErrorType("not found")),
		errors.SetLocalizer(i18n.NewLocalizer(bundle, "ru")),
		errors.SetTranslateContext(&errors.TranslateContext{
			TemplateData: map[string]interface{}{
				"Name":        "John Snow",
				"PluralCount": 5,
			},
			PluralCount: 5,
		}),
	)

	fmt.Printf("%v\n", err)
}
