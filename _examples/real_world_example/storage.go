package main

import (
	"context"
	"time"

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

	if boorand.Bool() {
		return nil, ErrRandomError
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

	if _, err := s.GetUserByID(ctx, id); err == nil {
		return ErrDuplicateKey
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

	result := s.db.Create(&usr)

	if boorand.Bool() {
		return ErrRandomError
	}

	return result.Error
}

func (s *Storage) DeleteUserByID(ctx context.Context, id string) error {
	result := s.db.Delete(&UserModel{ID: id})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if boorand.Bool() {
		return ErrRandomError
	}

	return result.Error
}
