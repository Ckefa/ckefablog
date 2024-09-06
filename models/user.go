package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    string `json:"id"`
	Email string `json:"email"`
}

func NewUser(email string) *User {
	return &User{
		ID:    uuid.New().String(),
		Email: email,
	}
}
