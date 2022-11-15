package domain

import (
	"errors"
	"fmt"
	"net/mail"

	"golang.org/x/exp/slices"
)

type User struct {
	ID     uint32 `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Gender Gender `json:"gender"`
	Age    uint8  `json:"age"`
}

func (u User) Validate() error {
	if u.Name == "" {
		return errors.New("empty name")
	}

	if !slices.Contains(Genders, u.Gender) {
		return fmt.Errorf("unknown gender %q, expected one of %q", u.Gender, Genders)
	}

	if u.Age < 18 {
		return errors.New("age must be equal or above 18")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	return nil
}
