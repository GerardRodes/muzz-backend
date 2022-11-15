package domain

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepoFilter struct {
	ExcludeUsers         []uint32
	GenderOpposedToUsers []uint32
	NotSwipedByUsers     []uint32
}

type UserRepo interface {
	Create(ctx context.Context, user User, passwordHash []byte) (id uint32, err error)
	List(ctx context.Context, filter UserRepoFilter) ([]User, error)
}

type UserSvc struct {
	r UserRepo
}

func NewUserSvc(r UserRepo) UserSvc {
	return UserSvc{r}
}

// CreateRandom creates a user with random values
func (s UserSvc) CreateRandom(ctx context.Context) (User, error) {
	name := fmt.Sprintf("random-%d", time.Now().UnixMilli())
	user := User{
		Email:  fmt.Sprintf("%s@muzz.com", name),
		Name:   name,
		Gender: Genders[rand.Intn(len(Genders))],
		Age:    uint8(rand.Intn(70) + 18),
	}

	var err error
	if user.ID, err = s.Create(ctx, user, name); err != nil {
		return User{}, fmt.Errorf("creating user: %w", err)
	}

	return user, nil
}

func (s UserSvc) Create(ctx context.Context, user User, password string) (uint32, error) {
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

	id, err := s.r.Create(ctx, user, passwordHash)
	if err != nil {
		return 0, fmt.Errorf("persisting user: %w", err)
	}

	return id, nil
}

// ListPotentialMatches lists profiles which:
//   - are not the user itself
//   - are from the opposed gender
//   - have not been already swiped by the user
func (s UserSvc) ListPotentialMatches(ctx context.Context, userID uint32) ([]User, error) {
	users, err := s.r.List(ctx, UserRepoFilter{
		ExcludeUsers:         []uint32{userID},
		GenderOpposedToUsers: []uint32{userID},
		NotSwipedByUsers:     []uint32{userID},
	})

	if err != nil {
		return nil, fmt.Errorf("list potential matches: %w", err)
	}

	return users, nil
}