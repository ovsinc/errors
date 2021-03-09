package errors

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/stretchr/testify/require"
)

func BenchmarkTranslateMsg(b *testing.B) {
	var (
		unreadEmailCount = 5
		name             = "John Snow"
	)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("./testdata/active.ru.toml")

	localizer := i18n.NewLocalizer(bundle, "ru")

	ErrEmailsUnreadMsg := TranslateMessage{
		ID: "ErrEmailsUnreadMsg",
		TemplateData: map[string]interface{}{
			"Name":        name,
			"PluralCount": unreadEmailCount,
		},
		PluralCount: unreadEmailCount,
		Localizer:   localizer,
	}

	e1 := New(
		"fallback message",
		SetErrorType(NewErrorType("not found")),
		AddTranslatedMessage("ru", &ErrEmailsUnreadMsg),
		SetLang("ru"),
	)

	require.Equal(b, e1.Error(), "[not found][ERROR] -- У John Snow имеется 5 непрочитанных сообщений.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e1.Error()
	}
}
