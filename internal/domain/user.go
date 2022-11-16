package domain

import (
	"errors"
	"fmt"
	"net/mail"
)

type User struct {
	ID       uint32 `json:"id"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name"`
	Gender   Gender `json:"gender"`
	Age      uint8  `json:"age"`
	Location Point  `json:"-"`
}

func (u User) Validate() error {
	if u.Name == "" {
		return errors.New("empty name")
	}

	if err := u.Gender.Validate(); err != nil {
		return err
	}

	if u.Age < 18 {
		return errors.New("age must be equal or above 18")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	return nil
}
