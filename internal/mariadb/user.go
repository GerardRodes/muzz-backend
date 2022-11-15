package mariadb

import (
	"context"
	"database/sql"

	"github.com/GerardRodes/muzz-backend/internal/domain"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) userRepo {
	return userRepo{db}
}

const userRepoCreateQuery = `
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

func (r userRepo) Create(ctx context.Context, user domain.User, passwordHash []byte) (id uint32, err error) {
	row := r.db.QueryRowContext(ctx, userRepoCreateQuery, user.Email, passwordHash, user.Name, user.Gender, user.Age)
	err = row.Scan(&id)
	return
}

func (r userRepo) List(ctx context.Context, filter domain.UserRepoFilter) ([]domain.User, error) {
	return nil, nil
}
