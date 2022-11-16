package mariadb

import (
	"context"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return Repo{db}
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
