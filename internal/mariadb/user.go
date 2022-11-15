package mariadb

import (
	"context"
	"database/sql"

	"github.com/GerardRodes/muzz-backend/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{db}
}

const userRepoCreate = `
INSERT INTO users (
	email,
	password,
	name,
	gender,
	age
) VALUES (
	?, ?, ?, ?, ?
)
RETURNING id
`

func (r UserRepo) Create(ctx context.Context, user domain.User) (id uint32, err error) {
	return 0, nil
}

func (r UserRepo) List(ctx context.Context, filter domain.UserRepoFilter) ([]domain.User, error) {
	return nil, nil
}
