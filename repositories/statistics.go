package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatRepo struct {
	db *pgxpool.Pool
}

func NewStatRepo(db *pgxpool.Pool) *StatRepo {
	return &StatRepo{
		db: db,
	}
}

func (r *StatRepo) CountDiary(
	ctx context.Context,
	user_id int64,
) (int64, error) {
	var total int64

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM diaries
		WHERE user_id = $1
		`, user_id,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, err
}

func (r *StatRepo) CountDiaryByMood(
	ctx context.Context,
	user_id int64,
) (map[int64]int64, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT 
		mood_id,
		COUNT(*)
		FROM diaries
		WHERE user_id = $1
		GROUP BY mood_id
		ORDER BY mood_id
		`, user_id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(map[int64]int64)

	for rows.Next() {
		var mood_id int64
		var total int64

		err := rows.Scan(
			&mood_id,
			&total,
		)

		if err != nil {
			return nil, err
		}

		result[mood_id] = total
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (r *StatRepo) CountDiaryByMonth(
	ctx context.Context,
	user_id int64,
	year int,
) (map[int]int64, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
		EXTRACT(MONTH FROM date)::int AS month,
		COUNT (*)
		FROM diaries
		WHERE user_id = $1
		AND EXTRACT(YEAR FROM date)::int = $2
		GROUP BY month
		ORDER BY month
		`, user_id, year,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(map[int]int64)

	for rows.Next() {
		var month int
		var total int64

		err := rows.Scan(
			&month,
			&total,
		)

		if err != nil {
			return nil, err
		}

		result[month] = total
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}
