package mariadb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GerardRodes/muzz-backend/internal/domain"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) userRepo {
	return userRepo{db}
}

const userCreateQuery = `
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
	row := r.db.QueryRowContext(ctx, userCreateQuery, user.Email, passwordHash, user.Name, user.Gender, user.Age)
	err = row.Scan(&id)
	return
}

const userGetQuery = `
SELECT
	email,
	name,
	gender,
	age
FROM users
WHERE id != ?
`

func (r userRepo) Get(ctx context.Context, userID uint32) (user domain.User, err error) {
	row := r.db.QueryRowContext(ctx, userGetQuery, userID)
	user.ID = userID
	err = row.Scan(&user.Email, &user.Name, &user.Gender, &user.Age)
	return
}

const userListQuery = `
SELECT
	id,
	name,
	gender,
	age
FROM users
WHERE id != ?
	AND gender != ?
`

func (r userRepo) ListPotentialMatches(ctx context.Context, user domain.User) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, userListQuery, user.ID, user.Gender)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Gender, &user.Age); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
