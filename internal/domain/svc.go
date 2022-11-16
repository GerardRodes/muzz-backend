package domain

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrWrongPassword = errors.New("wrong password")

type Repo interface {
	GetUser(ctx context.Context, userID uint32) (User, error)
	GetUserIDAndPasswordByEmail(ctx context.Context, email string) (userID uint32, passHash []byte, err error)
	CreateUser(ctx context.Context, user User, passwordHash []byte) (uint32, error)
	ListPotentialMatches(ctx context.Context, user User) ([]User, error)
	Swipe(ctx context.Context, userID, profileID uint32, preference bool) error
	BothLiked(ctx context.Context, userID1, userID2 uint32) (bool, error)
	CreateMatch(ctx context.Context, userID1, userID2 uint32) (uint64, error)
}

type SessionStorage interface {
	Create(ctx context.Context, userID uint32) (sessionID string, err error)
	Load(ctx context.Context, sessionID string) (userID uint32, err error)
}

type Service struct {
	r  Repo
	ss SessionStorage
}

func NewService(r Repo, ss SessionStorage) Service {
	return Service{r, ss}
}

// ListPotentialMatches lists profiles which:
//   - are not the user itself
//   - are from the opposed gender
//   - have not been already swiped by the user
func (s Service) ListPotentialMatches(ctx context.Context, userID uint32) ([]User, error) {
	user, err := s.r.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	users, err := s.r.ListPotentialMatches(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}

	return users, nil
}

func (s Service) Swipe(ctx context.Context, userID, profileID uint32, preference bool) (uint64, error) {
	if err := s.r.Swipe(ctx, userID, profileID, preference); err != nil {
		return 0, fmt.Errorf("cannot add swipe: %w", err)
	}

	liked, err := s.r.BothLiked(ctx, userID, profileID)
	if err != nil {
		return 0, fmt.Errorf("cannot check if both liked: %w", err)
	}

	if !liked {
		return 0, nil
	}

	matchID, err := s.r.CreateMatch(ctx, userID, profileID)
	if err != nil {
		return 0, fmt.Errorf("cannot create match: %w", err)
	}

	return matchID, nil
}

func (s Service) Login(ctx context.Context, email, password string) (string, error) {
	userID, passHash, err := s.r.GetUserIDAndPasswordByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrWrongPassword
		}
		return "", fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(passHash, []byte(password)); err != nil {
		return "", ErrWrongPassword
	}

	sessionID, err := s.ss.Create(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return sessionID, nil
}

func (s Service) LoadSession(ctx context.Context, sessionID string) (uint32, error) {
	return s.ss.Load(ctx, sessionID)
}
