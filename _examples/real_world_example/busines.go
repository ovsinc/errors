package main

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"

	"github.com/ovsinc/errors"
)

var (
	ErrUserValidation = errors.NewWith(
		errors.SetMsg(EValidationMsg.Other),
		errors.SetID(EValidationMsg.ID),
		errors.SetErrorType(EValidation.String()),
	)

	ErrUserEmpty = errors.NewWith(
		errors.SetID(EEmptyMsg.ID),
		errors.SetMsg(EEmptyMsg.Other),
		errors.SetErrorType(EEmpty.String()),
	)

	ErrDBNotFound = errors.NewWith(
		errors.SetMsg(ENotFoundMsg.Other),
		errors.SetID(ENotFoundMsg.ID),
		errors.SetErrorType(ENotFound.String()),
	)
	ErrDBDuplicate = errors.NewWith(
		errors.SetMsg(EDuplicateMsg.Other),
		errors.SetID(EDuplicateMsg.ID),
		errors.SetErrorType(EDuplicate.String()),
	)
	ErrDBInternal = errors.NewWith(
		errors.SetMsg(EInternalMsg.Other),
		errors.SetID(EInternalMsg.ID),
		errors.SetErrorType(EInternal.String()),
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
		validation.Field(&u.FirstName, validation.Required, validation.Length(5, 20)),
		validation.Field(&u.LastName, validation.Required, validation.Length(5, 20)),
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
			errors.Wrap(ErrDBNotFound.WithOptions(errors.SetOperation("storage.GetUserByID")), err)

	default:
		return nil,
			errors.Wrap(ErrDBInternal.WithOptions(errors.SetOperation("storage.GetUserByID")), err)
	}
}

func (uc *UserUC) NewUser(ctx context.Context, u *User) error {
	if u == nil {
		return ErrUserEmpty.WithOptions(errors.SetOperation("uc.NewUser"))
	}

	if err := u.Validate(); err != nil {
		return errors.Wrap(ErrUserValidation.WithOptions(errors.SetOperation("uc.NewUser")), err)
	}

	if err := uc.userStore.CreateUser(ctx, u); err != nil {
		return errors.Wrap(ErrDBInternal.WithOptions(errors.SetOperation("uc.NewUser")), err)
	}

	return nil
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
		return errors.Wrap(ErrDBNotFound.WithOptions(errors.SetOperation("storage.DeleteUserByID")), err)

	default:
		return errors.Wrap(ErrDBInternal.WithOptions(errors.SetOperation("storage.DeleteUserByID")), err)
	}
}
