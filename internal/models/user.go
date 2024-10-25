package models

import "time"


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