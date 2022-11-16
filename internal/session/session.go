package session

import (
	"context"
	"time"

	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type SessionStorage struct {
	rdb        *redis.Client
	expiration time.Duration
}

type Config struct {
	RedisClient *redis.Client
	Expiration  time.Duration
}

func NewSessionStorage(c Config) SessionStorage {
	return SessionStorage{
		expiration: c.Expiration,
		rdb:        c.RedisClient,
	}
}

func (s SessionStorage) Create(ctx context.Context, userID uint32) (sessionID string, err error) {
	sessionID = uuid.New().String()
	return sessionID, s.rdb.Set(ctx, sessionID, userID, s.expiration).Err()
}

func (s SessionStorage) Load(ctx context.Context, sessionID string) (userID uint32, err error) {
	v, err := s.rdb.Get(ctx, sessionID).Uint64()
	if err == redis.Nil {
		return 0, domain.ErrNotFound
	}
	return uint32(v), err
}
