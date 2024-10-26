package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)


type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetID returns the user ID.
func (u User) GetID() int {
	return u.ID
}

// GetName returns the user name.
func (u User) GetEmail() string {
	return u.Email
}

func (User) TableName() string {
    return "users"
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest User

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r, validation.Field(&r.Email, is.Email, validation.Required), validation.Field(&r.Password, validation.Required, validation.Length(6,100)))
}

func (r User) Validate() error {
    return validation.ValidateStruct(&r, validation.Field(&r.Name, validation.Required), validation.Field(&r.Email, is.Email, validation.Required), validation.Field(&r.Password, validation.Required, validation.Length(6,100)))
}