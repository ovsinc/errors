package main

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"

	"github.com/ovsinc/errors"
)

var (
	ErrUserValidation = errors.ValidationErrWith(
		errors.SetMsg(EValidationMsg.Other),
		errors.SetID(EValidationMsg.ID),
	)

	ErrUserEmpty = errors.EmptyErrWith(
		errors.SetID(EEmptyMsg.ID),
		errors.SetMsg(EEmptyMsg.Other),
	)

	ErrDBNotFound = errors.NotFoundErrWith(
		errors.SetMsg(ENotFoundMsg.Other),
		errors.SetID(ENotFoundMsg.ID),
	)
	ErrDBDuplicate = errors.DuplicateErrWith(
		errors.SetMsg(EDuplicateMsg.Other),
		errors.SetID(EDuplicateMsg.ID),
	)
	ErrDBInternal = errors.IternalErrWith(
		errors.SetMsg(EInternalMsg.Other),
		errors.SetID(EInternalMsg.ID),
	)

	ErrUnknown = errors.NewWith(
		errors.SetMsg(EUnknownMsg.Other),
		errors.SetID(EUnknownMsg.ID),
		errors.SetErrorType(errors.Unknown),
	)
)

type User struct {
	ID        string `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Email     string `json:"email" gorm:"index;not null"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 20)),
		validation.Field(&u.LastName, validation.Required, validation.Length(1, 20)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.ID, validation.Required, is.UUID),
	)
}

//

func NewUserUC(userStore *Storage) *UserUC {
	return &UserUC{userStore: userStore}
}

type UserUC struct {
	userStore *Storage
}

func (uc *UserUC) GetUser(ctx context.Context, id string) (*User, error) {
	if err := validation.Validate(&id, validation.Required, is.UUID); err != nil {
		return nil,
			errors.Wrap(ErrUserValidation.WithOptions(errors.SetOperation("uc.GetUser")), err)
	}

	usr, err := uc.userStore.GetUserByID(ctx, id)
	switch {
	case err == nil:
		return usr, nil

	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil,
			errors.Wrap(ErrDBNotFound.WithOptions(errors.SetOperation("uc.GetUser")), err)

	case errors.Is(err, ErrRandomError):
		return nil,
			errors.Wrap(ErrUnknown.WithOptions(errors.SetOperation("uc.GetUser")), err)

	default:
		return nil,
			errors.Wrap(ErrDBInternal.WithOptions(errors.SetOperation("uc.GetUser")), err)
	}
}

func (uc *UserUC) NewUser(ctx context.Context, u *User) error {
	if u == nil {
		return ErrUserEmpty.WithOptions(errors.SetOperation("uc.NewUser"))
	}

	if err := u.Validate(); err != nil {
		return errors.Wrap(ErrUserValidation.WithOptions(errors.SetOperation("uc.NewUser")), err)
	}

	err := uc.userStore.CreateUser(ctx, u)
	switch {
	case err == nil:
		return nil

	case errors.Is(err, ErrDuplicateKey):
		return errors.Wrap(ErrDBDuplicate.WithOptions(errors.SetOperation("uc.NewUser")), err)

	case errors.Is(err, ErrRandomError):
		return errors.Wrap(ErrUnknown.WithOptions(errors.SetOperation("uc.NewUser")), err)

	default:
		return errors.Wrap(ErrDBInternal.WithOptions(errors.SetOperation("uc.GetUser")), err)
	}
}

func (uc *UserUC) DeleteUser(ctx context.Context, id string) error {
	if err := validation.Validate(&id, validation.Required, is.UUID); err != nil {
		return errors.Wrap(ErrUserValidation.WithOptions(errors.SetOperation("uc.DeleteUser")), err)
	}

	err := uc.userStore.DeleteUserByID(ctx, id)
	switch {
	case err == nil:
		return nil

	case errors.Is(err, gorm.ErrRecordNotFound):
		return errors.Wrap(ErrDBNotFound.WithOptions(errors.SetOperation("uc.DeleteUser")), err)

	case errors.Is(err, ErrRandomError):
		return errors.Wrap(ErrUnknown.WithOptions(errors.SetOperation("uc.DeleteUser")), err)

	default:
		return errors.Wrap(ErrDBInternal.WithOptions(errors.SetOperation("uc.DeleteUser")), err)
	}
}
