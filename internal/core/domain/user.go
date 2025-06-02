package domain

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty" bson:"id" `                      // auto-generated
	Name      string    `json:"name" bson:"name" validate:"required"`         // string
	Email     string    `json:"email" bson:"email" validate:"required,email"` // unique
	Password  string    `json:"password" bson:"password" validate:"required"` // hashed
	CreatedAt time.Time `json:"created_at" bson:"created_at"`                 // timestamp
}

func (u *User) ValidateEmailAndName() error {
	if u.Email == "" || u.Name == "" {
		return errors.New("name and email cannot be empty")
	}
	return nil
}
