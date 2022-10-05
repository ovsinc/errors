package main

import (
	"embed"

	"github.com/BurntSushi/toml"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/ovsinc/errors"
	"golang.org/x/text/language"

	"io/fs"
)

const (
	ErrBadContentID     = "ErrBadContent"
	ErrUserValidationID = "ErrUserValidation"
	ErrUserEmptyID      = "ErrUserEmpty"

	ErrDBNotFoundID  = "ErrDBNotFound"
	ErrDBDuplicateID = "ErrDBDuplicate"
	ErrDBInternalID  = "ErrDBInternal"

	ErrUnknownID = "ErrUnknown"
)

var (
	EEmptyMsg = &i18n.Message{
		ID:          ErrUserEmptyID,
		Other:       "Empty entity.",
		Description: "Входные данные пусты или nil",
	}
	EValidationMsg = &i18n.Message{
		ID:          ErrUserValidationID,
		Other:       "Validation error.",
		Description: "Провал валидации",
	}
	ENotFoundMsg = &i18n.Message{
		ID:          ErrDBNotFoundID,
		Other:       "Entity not found.",
		Description: "Запись в базе не найдена",
	}
	EDuplicateMsg = &i18n.Message{
		ID:          ErrDBDuplicateID,
		Other:       "Duplicate entity.",
		Description: "Такая запись в базе уже естьб дубликат",
	}
	EInternalMsg = &i18n.Message{
		ID:          ErrDBInternalID,
		Other:       "Internal error.",
		Description: "Любая внутренняя ошибка",
	}
	EBadContentMsg = &i18n.Message{
		ID:          ErrBadContentID,
		Other:       "Incorrect input data.",
		Description: "Некорректные входные данные",
	}
	EUnknownMsg = &i18n.Message{
		ID:          ErrUnknownID,
		Other:       "Unknown error.",
		Description: "Внезапная ошибка",
	}
)

//go:embed translate/active.*.toml
var LocaleFS embed.FS

func NewTranslater() *translater {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	paths, err := fs.Glob(LocaleFS, "translate/active.*.toml")
	if err != nil {
		panic(err)
	}
	for _, path := range paths {
		_, err := bundle.LoadMessageFileFS(LocaleFS, path)
		if err != nil {
			panic(err)
		}
	}

	return &translater{
		bundle: bundle,
	}
}

type translater struct {
	bundle *i18n.Bundle
}

func (t *translater) TranslateError(lang string, err error) string {
	localizer := i18n.NewLocalizer(t.bundle, lang)
	s, _ := errors.Translate(err, localizer, nil)
	return s
}
