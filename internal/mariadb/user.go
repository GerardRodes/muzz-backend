package mariadb

import (
	"context"

	"github.com/GerardRodes/muzz-backend/internal/domain"
)

const createUserQuery = `
INSERT INTO users (
	email,
	password,
	name,
	gender,
	age
) VALUES (
	?, ?, ?, ?, ?
) RETURNING id
`

func (r Repo) CreateUser(ctx context.Context, user domain.User, passwordHash []byte) (id uint32, err error) {
	row := r.db.QueryRowContext(ctx, createUserQuery, user.Email, passwordHash, user.Name, user.Gender, user.Age)
	err = row.Scan(&id)
	return
}

const getUserQuery = `
SELECT
	email,
	name,
	gender,
	age
FROM users
WHERE id = ?
`

func (r Repo) GetUser(ctx context.Context, userID uint32) (user domain.User, err error) {
	row := r.db.QueryRowContext(ctx, getUserQuery, userID)
	user.ID = userID
	err = row.Scan(&user.Email, &user.Name, &user.Gender, &user.Age)
	return
}
