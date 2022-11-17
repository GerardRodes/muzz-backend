package domain

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// userRandomCounter allows to create multiple accounts on same timestamp
// without email colision, so we can execute comands like this:
//
//	❯ seq 1 1000 | xargs -P1000 -n1 curl -XPOST http://localhost:8080/user/create
//
// and safely create 1000 users without CPU locks
var userRandomCounter atomic.Uint32

// CreateRandomUser creates a user with random values
func (s Service) CreateRandomUser(ctx context.Context) (User, string, error) {
	name := fmt.Sprintf("random-%d-%d", time.Now().UnixMilli(), userRandomCounter.Add(1))
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
