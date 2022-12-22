package main

import (
	"time"

	"github.com/ovsinc/errors"
	"gorm.io/gorm"
)

var ErrDuplicateKey = errors.New("duplacate record")

type UserModel struct {
	gorm.Model

	ID        string `gorm:"primaryKey"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"index;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
