package model

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

var (
	ErrEmailEmpty    = errors.New("email cannot be empty")
	ErrPasswordEmpty = errors.New("password cannot be empty")
)

type User struct {
	UUID     uuid.UUID `json:"user_id" bson:"_id"`
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"password" bson:"password"`
}

type Input struct {
	User
}

func (i *Input) Validate() error {
	switch {
	case i.Email == "":
		return ErrEmailEmpty
	case i.Password == "":
		return ErrPasswordEmpty
	default:
		return nil
	}
}
func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
