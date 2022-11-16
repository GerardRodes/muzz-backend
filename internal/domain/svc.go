package domain

import (
	"context"
	"fmt"
)

type Repo interface {
	GetUser(ctx context.Context, userID uint32) (User, error)
	CreateUser(ctx context.Context, user User, passwordHash []byte) (uint32, error)
	ListPotentialMatches(ctx context.Context, user User) ([]User, error)
	Swipe(ctx context.Context, userID, profileID uint32, preference bool) error
	BothLiked(ctx context.Context, userID1, userID2 uint32) (bool, error)
	CreateMatch(ctx context.Context, userID1, userID2 uint32) (uint64, error)
}

type Service struct {
	r Repo
}

func NewService(r Repo) Service {
	return Service{r}
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
