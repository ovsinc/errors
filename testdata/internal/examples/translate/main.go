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
		errors.SetErrorType(errors.NewErrorType("not found")),
		errors.AddTranslatedMessage("ru", &errors.TranslateMessage{
			TemplateData: map[string]interface{}{
				"Name":        "John Snow",
				"PluralCount": 5,
			},
			PluralCount: 5,
			Localizer:   i18n.NewLocalizer(bundle, "ru"),
			ID:          "ErrEmailsUnreadMsg",
		}),
		errors.SetLang("ru"),
	)

	fmt.Printf("%v\n", err)
}
