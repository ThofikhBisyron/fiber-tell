package repositories

import (
	"context"
	"fmt"
	"tell-be/models"
	"time"

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

func (r *StatRepo) CountDiaryByPeriod(
	ctx context.Context,
	user_id int64,
	start_date time.Time,
	end_date time.Time,
) (int64, error) {
	var total int64

	err := r.db.QueryRow(
		ctx,
		`
        SELECT COUNT(*)
        FROM diaries
        WHERE user_id = $1
        AND date >= $2
        AND date < $3
        `,
		user_id,
		start_date,
		end_date,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *StatRepo) CountDiaryByMoodAndPeriod(
	ctx context.Context,
	user_id int64,
	start_date time.Time,
	end_date time.Time,
) (map[int64]int64, error) {
	rows, err := r.db.Query(
		ctx,
		`
        SELECT mood_id, COUNT(*)
        FROM diaries
        WHERE user_id = $1
        AND date >= $2
        AND date < $3
        GROUP BY mood_id
        ORDER BY mood_id
        `,
		user_id,
		start_date,
		end_date,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(map[int64]int64)

	for rows.Next() {
		var mood_id int64
		var total int64

		if err := rows.Scan(&mood_id, &total); err != nil {
			return nil, err
		}

		result[mood_id] = total
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *StatRepo) CountDiaryActivity(
	ctx context.Context,
	user_id int64,
	start_date time.Time,
	end_date time.Time,
	period string,
) ([]models.DiaryActivity, error) {

	rows, err := r.db.Query(
		ctx,
		`
        SELECT
            date,
            COUNT(*)
        FROM diaries
        WHERE user_id = $1
        AND date >= $2
        AND date < $3
        GROUP BY date
        ORDER BY date
        `,
		user_id,
		start_date,
		end_date,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int64)

	for rows.Next() {
		var date time.Time
		var total int64

		if err := rows.Scan(&date, &total); err != nil {
			return nil, err
		}

		counts[date.Format("2006-01-02")] = total
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]models.DiaryActivity, 0)

	switch period {
	case "weekly":
		for i := 0; i < 7; i++ {
			date := start_date.AddDate(0, 0, i)
			key := date.Format("2006-01-02")

			result = append(result, models.DiaryActivity{
				Label: date.Format("Mon"),
				Total: counts[key],
			})
		}

	case "monthly":
		for date := start_date; date.Before(end_date); date = date.AddDate(0, 0, 1) {
			key := date.Format("2006-01-02")

			result = append(result, models.DiaryActivity{
				Label: date.Format("02"),
				Total: counts[key],
			})
		}

	case "yearly":
		monthlyCounts := make(map[string]int64)

		for key, total := range counts {
			date, err := time.Parse("2006-01-02", key)
			if err != nil {
				return nil, err
			}

			monthKey := date.Format("2006-01")
			monthlyCounts[monthKey] += total
		}

		for i := 0; i < 12; i++ {
			date := start_date.AddDate(0, i, 0)
			key := date.Format("2006-01")

			result = append(result, models.DiaryActivity{
				Label: date.Format("Jan"),
				Total: monthlyCounts[key],
			})
		}

	default:
		return nil, fmt.Errorf("invalid period")
	}

	return result, nil
}
