package mariadb

import (
	"context"
	"fmt"

	"github.com/GerardRodes/muzz-backend/internal/domain"
)

func (r Repo) ListPotentialMatches(
	ctx context.Context,
	user domain.User,
	filter domain.ListPotentialMatchesFilter,
) ([]domain.ListPotentialMatchesResult, error) {
	var query string
	args := []any{user.Location.Lng, user.Location.Lat, user.ID, user.ID}

	if filter.OrderByLikes {
		query += `
			WITH attractiveness AS (
				SELECT u.id AS "user_id", 2*SUM(s.preference) - COUNT(s.user_id) as "score"
				FROM users u
				LEFT JOIN swipes s ON s.profile_id = u.id
				GROUP BY u.id
			)
		`
	}

	query += `
		SELECT
			u.id,
			u.name,
			u.gender,
			u.age,
			ST_DISTANCE_SPHERE(Point(?, ?), u.location) as "distanceFromMe"
		FROM users u
	`

	if filter.OrderByLikes {
		query += `
			LEFT JOIN attractiveness a ON a.user_id = u.id
		`
	}

	query += `
		WHERE u.id != ?
			AND u.id NOT IN (
				SELECT s.profile_id FROM swipes s WHERE s.user_id = ?
			)
	`

	if filter.AgeMin != 0 {
		query += `
			AND u.age >= ?
		`
		args = append(args, filter.AgeMin)
	}
	if filter.AgeMax != 0 {
		query += `
			AND u.age <= ?
		`
		args = append(args, filter.AgeMax)
	}
	if filter.Gender != "" {
		query += `
			AND u.gender = ?
		`
		args = append(args, filter.Gender)
	}

	if filter.OrderByLikes {
		query += `
			ORDER BY a.score DESC
		`
	} else {
		query += `
			ORDER BY distanceFromMe ASC
		`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var results []domain.ListPotentialMatchesResult
	for rows.Next() {
		var result domain.ListPotentialMatchesResult
		if err := rows.Scan(&result.ID, &result.Name, &result.Gender, &result.Age, &result.DistanceFromMe); err != nil {
			return nil, fmt.Errorf("scan result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}
