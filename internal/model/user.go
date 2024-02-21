package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUUIDNil          = errors.New("UUID cannot be nil")
	ErrNameEmpty        = errors.New("name cannot be empty")
	ErrEmailEmpty       = errors.New("email cannot be empty")
	ErrPasswordEmpty    = errors.New("password cannot be empty")
	ErrRegisteredAtZero = errors.New("registeredAt cannot be zero")
)

type User struct {
	UUID         uuid.UUID `json:"user_id" bson:"_id"`
	Name         string    `json:"name" bson:"name"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"password" bson:"password"`
	RegisteredAt time.Time `json:"registeredAt" bson:"registeredAt"`
}

type SignUpInput struct {
	UUID     uuid.UUID `json:"user_id" bson:"_id"`
	Name     string    `json:"name" bson:"name"`
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"password" bson:"password"`
}

func (u *User) Validate() error {
	switch {
	case u.UUID == uuid.Nil:
		return ErrUUIDNil
	case u.Name == "":
		return ErrNameEmpty
	case u.Email == "":
		return ErrEmailEmpty
	case u.Password == "":
		return ErrPasswordEmpty
	case u.RegisteredAt.IsZero():
		return ErrRegisteredAtZero
	default:
		return nil
	}
}
