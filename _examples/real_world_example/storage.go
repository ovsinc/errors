package main

import (
	"context"
	"time"

	"github.com/ovsinc/errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewStorage() *Storage {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&UserModel{}); err != nil {
		panic(err)
	}

	return &Storage{db: db}
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) GetUserByID(_ context.Context, id string) (*User, error) {
	usr := new(UserModel)

	if result := s.db.First(usr); result.Error != nil {
		return nil, result.Error
	}

	return &User{
		ID:        usr.ID,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.LastName,
	}, nil
}

func (s *Storage) CreateUser(ctx context.Context, u *User) error {
	id := u.ID

	//sqite драйвер не поддерживает проверку уникальности, нет ошибки ErrDuplicateKey
	if _, err := s.GetUserByID(ctx, id); err == nil {
		return ErrDBDuplicate.WithOptions(errors.SetOperation("CreateUser"))
	}

	now := time.Now()

	usr := UserModel{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		UpdatedAt: now,
		CreatedAt: now,
	}

	if result := s.db.Create(&usr); result.Error != nil {
		return errors.Wrap(
			ErrDBInternal.WithOptions(errors.SetOperation("CreateUser")),
			result.Error)
	}

	return nil
}

func (s *Storage) DeleteUserByID(ctx context.Context, id string) error {
	result := s.db.Delete(&UserModel{}, id)

	switch {
	case result.Error == nil:
		return nil

	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return errors.Wrap(
			ErrDBNotFound.WithOptions(errors.SetOperation("DeleteUserByID")),
			result.Error)

	default:
		return errors.Wrap(
			ErrDBInternal.WithOptions(errors.SetOperation("DeleteUserByID")),
			result.Error)
	}
}
