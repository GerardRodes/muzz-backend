package mariadb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GerardRodes/muzz-backend/internal/domain"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return Repo{db}
}

const listPotentialMatchesQuery = `
SELECT
	u.id,
	u.name,
	u.gender,
	u.age
FROM users u
WHERE u.id != ?
	AND u.id NOT IN (
		SELECT s.profile_id FROM swipes s WHERE s.user_id = ?
	)
`
const listPotentialMatchesMinAgeFilter = `
	AND u.age >= ?
`
const listPotentialMatchesMaxAgeFilter = `
	AND u.age <= ?
`
const listPotentialMatchesGenderFilter = `
	AND u.gender = ?
`

func (r Repo) ListPotentialMatches(ctx context.Context, user domain.User, filter domain.ListPotentialMatchesFilter) ([]domain.User, error) {
	query := listPotentialMatchesQuery
	args := []any{user.ID, user.ID}
	if filter.AgeMin != 0 {
		query += listPotentialMatchesMinAgeFilter
		args = append(args, filter.AgeMin)
	}
	if filter.AgeMax != 0 {
		query += listPotentialMatchesMaxAgeFilter
		args = append(args, filter.AgeMax)
	}
	if filter.Gender != "" {
		query += listPotentialMatchesGenderFilter
		args = append(args, filter.Gender)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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

const swipeQuery = `
INSERT INTO swipes (
	user_id,
	profile_id,
	preference
) VALUES (
	?, ?, ?
) ON DUPLICATE KEY UPDATE user_id=user_id
`

func (r Repo) Swipe(ctx context.Context, userID, profileID uint32, preference bool) error {
	_, err := r.db.ExecContext(ctx, swipeQuery, userID, profileID, preference)
	return err
}

const bothLikedQuery = `
SELECT TRUE
FROM swipes
WHERE user_id IN (?, ?)
	AND profile_id IN (?, ?)
	AND preference IS TRUE
HAVING COUNT(*) > 1
`

func (r Repo) BothLiked(ctx context.Context, userID1, userID2 uint32) (bool, error) {
	row := r.db.QueryRowContext(ctx, bothLikedQuery, userID1, userID2, userID1, userID2)
	var liked bool
	err := row.Scan(&liked)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return liked, err
}

const createMatchQuery = `
INSERT INTO matches (
	user_id_low,
	user_id_high
) VALUES (
	?, ?
) ON DUPLICATE KEY UPDATE id=id
RETURNING id
`

func (r Repo) CreateMatch(ctx context.Context, userID1, userID2 uint32) (uint64, error) {
	low, high := userID1, userID2
	if userID1 > userID2 {
		low, high = userID2, userID1
	}
	row := r.db.QueryRowContext(ctx, createMatchQuery, low, high)
	var matchID uint64
	err := row.Scan(&matchID)
	return matchID, err
}
