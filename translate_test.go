package errors

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func BenchmarkTranslateMsg(b *testing.B) {
	DefaultMultierrFormatFunc = StringMultierrFormatFunc

	var (
		unreadEmailCount = 5
		name             = "John Snow"
	)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("./internal/examples/translate/testdata/active.ru.toml")

	localizer := i18n.NewLocalizer(bundle, "ru")

	ErrEmailsUnreadMsg := TranslateContext{
		TemplateData: map[string]interface{}{
			"Name":        name,
			"PluralCount": unreadEmailCount,
		},
		PluralCount: unreadEmailCount,
	}

	e1 := New(
		"fallback message",
		SetID("ErrEmailsUnreadMsg"),
		SetLocalizer(localizer),
		SetErrorType("not found"),
		SetTranslateContext(&ErrEmailsUnreadMsg),
	)

	require.Equal(b, e1.Error(), "(not found) -- У John Snow имеется 5 непрочитанных сообщений.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e1.Error()
	}
}
