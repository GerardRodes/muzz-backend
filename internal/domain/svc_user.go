package domain

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// CreateRandomUser creates a user with random values
func (s Service) CreateRandomUser(ctx context.Context) (User, string, error) {
	name := fmt.Sprintf("random-%d", time.Now().UnixMilli())
	user := User{
		Email:  fmt.Sprintf("%s@muzz.com", name),
		Name:   name,
		Gender: Genders[rand.Intn(len(Genders))],
		Age:    uint8(rand.Intn(70) + 18),
		Location: Point{
			Lat: float64(-89+rand.Intn(179)) + rand.Float64(),
			Lng: float64(-179+rand.Intn(359)) + rand.Float64(),
		},
	}

	var err error
	if user.ID, err = s.CreateUser(ctx, user, name); err != nil {
		return User{}, "", fmt.Errorf("creating user: %w", err)
	}

	return user, name, nil
}

func (s Service) CreateUser(ctx context.Context, user User, password string) (uint32, error) {
	if err := user.Validate(); err != nil {
		return 0, fmt.Errorf("validating user: %w", err)
	}

	if len(password) < 6 {
		return 0, errors.New("password must have at least 6 characters")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hashing password: %w", err)
	}

	id, err := s.r.CreateUser(ctx, user, passwordHash)
	if err != nil {
		return 0, fmt.Errorf("persisting user: %w", err)
	}

	return id, nil
}
