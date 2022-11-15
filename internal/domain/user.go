package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"net/mail"
	"time"

	"golang.org/x/exp/slices"
)

type User struct {
	ID     uint32
	Email  string
	Name   string
	Gender Gender
	Age    uint8
}

func (u User) Validate() error {
	if u.Name == "" {
		return errors.New("user validation: empty name")
	}

	if !slices.Contains(Genders, u.Gender) {
		return fmt.Errorf("user validation: unknown gender %q, expected one of %q", u.Gender, Genders)
	}

	if u.Age < 18 {
		return errors.New("user validation: age must be equal or above 18")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("user validation: invalid email: %w", err)
	}

	return nil
}

type UserRepo interface {
	Create(u User) (id uint32, err error)
}

type UserSvc struct {
	r UserRepo
}

func (s UserSvc) CreateRandom() (User, error) {
	name := fmt.Sprintf("random-%d", time.Now().UnixMilli())
	u := User{
		Email:  fmt.Sprintf("%s@muzz.com", name),
		Name:   name,
		Gender: Genders[rand.Intn(len(Genders))],
		Age:    uint8(rand.Intn(70) + 18),
	}

	var err error
	if u.ID, err = s.Create(u); err != nil {
		return User{}, fmt.Errorf("create random user: %w", err)
	}

	return u, nil
}

func (s UserSvc) Create(u User) (uint32, error) {
	if err := u.Validate(); err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	id, err := s.r.Create(u)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}
