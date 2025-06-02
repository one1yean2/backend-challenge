package domain

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id" bson:"id,omitempty"`       // auto-generated
	Name      string    `json:"name" omitempty bson:"name"`   // string
	Email     string    `json:"email" omitempty bson:"email"` // unique
	Password  string    `json:"password" bson:"password"`     // hashed
	CreatedAt time.Time `json:"created_at" bson:"created_at"` // timestamp
}

func (u *User) ValidateEmailAndName() error {
	if u.Email == "" || u.Name == "" {
		return errors.New("name and email cannot be empty")
	}
	return nil
}
