package repositories

import (
	"context"
	"tell-be/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DiaryRepo struct {
	db *pgxpool.Pool
}

func NewDiaryRepo(db *pgxpool.Pool) *DiaryRepo {
	return &DiaryRepo{
		db: db,
	}
}

func (r *DiaryRepo) CreateDiary(
	ctx context.Context,
	user_id int64,
	mood_id int64,
	title string,
	content string,
	date time.Time,
) (models.Diaries, error) {
	var diary models.Diaries

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO diaries (
		user_id,
		mood_id,
		title,
		content,
		date
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
		id,
		user_id,
		mood_id,
		title,
		date,
		created_at,
		updated_at
		`, user_id, mood_id, title, content, date,
	).Scan(
		&diary.Id,
		&diary.User_id,
		&diary.Mood_id,
		&diary.Title,
		&diary.Content,
		&diary.Date,
		&diary.Created_at,
		&diary.Updated_at,
	)

	if err != nil {
		return models.Diaries{}, err
	}

	return diary, err
}

func (r *DiaryRepo) FindRecentDiary(
	ctx context.Context,
	user_id int64,
) (models.Diaries, error) {
	var diary models.Diaries

	err := r.db.QueryRow(
		ctx,
		`SELECT
		id,
        user_id,
        mood_id,
        title,
        content,
        date,
        created_at,
        updated_at
		FROM diaries
		WHERE user_id = $1
		ORDER BY date DESC, id DESC
		LIMIT 1`, user_id,
	).Scan(
		&diary.Id,
		&diary.User_id,
		&diary.Mood_id,
		&diary.Title,
		&diary.Content,
		&diary.Date,
		&diary.Created_at,
		&diary.Updated_at,
	)

	if err != nil {
		return models.Diaries{}, err
	}

	return diary, nil
}

func (r *DiaryRepo) FindDiariesByMonth(
	ctx context.Context,
	user_id int64,
	year int,
	month int,
) ([]models.Diaries, error) {
	var diaries []models.Diaries

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
		id,
        user_id,
        mood_id,
        title,
        content,
        date,
        created_at,
        updated_at
		FROM diaries
		WHERE user_id = $1
		AND date >= make_date($2, $3, 1)
		AND date < make_date($2, $3, 1) + INTERVAL '1 month'
		ORDER BY date DESC, id DESC
		`, user_id, year, month,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var diary models.Diaries

		err := rows.Scan(
			&diary.Id,
			&diary.User_id,
			&diary.Mood_id,
			&diary.Title,
			&diary.Content,
			&diary.Date,
			&diary.Created_at,
			&diary.Updated_at,
		)

		if err != nil {
			return nil, err
		}

		diaries = append(diaries, diary)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return diaries, nil
}

func (r *DiaryRepo) FindDIariesByDate(
	ctx context.Context,
	user_id int64,
	date time.Time,
) ([]models.Diaries, error) {
	var diaries []models.Diaries

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
		id,
        user_id,
        mood_id,
        title,
        content,
        date,
        created_at,
        updated_at
		FROM diaries
		WHERE user_id = $1
		AND date = $2
		ORDER BY id DESC
		`, user_id, date,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var diary models.Diaries

		err := rows.Scan(
			&diary.Id,
			&diary.User_id,
			&diary.Mood_id,
			&diary.Title,
			&diary.Content,
			&diary.Date,
			&diary.Created_at,
			&diary.Updated_at,
		)

		if err != nil {
			return nil, err
		}

		diaries = append(diaries, diary)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return diaries, err
}

func (r *DiaryRepo) FindDiaryById(
	ctx context.Context,
	user_id int64,
	diary_id int64,
) (models.Diaries, error) {
	var diary models.Diaries

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
		id,
        user_id,
        mood_id,
        title,
        content,
        date,
        created_at,
        updated_at
		FROM diaries
		WHERE id = $1
		AND user_id =$2
		`, diary_id, user_id,
	).Scan(
		&diary.Id,
		&diary.User_id,
		&diary.Mood_id,
		&diary.Title,
		&diary.Content,
		&diary.Date,
		&diary.Created_at,
		&diary.Updated_at,
	)

	if err != nil {
		return models.Diaries{}, err
	}

	return diary, nil
}

func (r *DiaryRepo) UpdateDiaryById(
	ctx context.Context,
	diary_id int64,
	user_id int64,
	mood_id int64,
	title string,
	content string,
	date time.Time,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE diaries
		SET
		mood_id = $1
		title = $2
		content = $3
		date = $4
		updated_at = NOW()
		WHERE id = $5
		AND user_id = $6
		`, mood_id, title, content, date, diary_id, user_id,
	)

	return err
}

func (r *DiaryRepo) DeleteDiaryById(
	ctx context.Context,
	diary_id int64,
	user_id int64,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM diaries 
		WHERE id = $1
		AND user_id = $2
		`, diary_id, user_id,
	)

	return err
}
