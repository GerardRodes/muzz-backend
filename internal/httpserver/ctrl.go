package httpserver

import (
	"context"

	"github.com/GerardRodes/muzz-backend/internal/domain"
)

type Service interface {
	CreateRandomUser(ctx context.Context) (domain.User, string, error)
	Swipe(ctx context.Context, userID, profileID uint32, preference bool) (matchID uint64, err error)
	ListPotentialMatches(ctx context.Context, userID uint32, filter domain.ListPotentialMatchesFilter) ([]domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	LoadSession(ctx context.Context, sessionID string) (uint32, error)
}

type Controller struct {
	svc Service
}
