package repositories

import (
	"context"
	"tell-be/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingRepo struct {
	db *pgxpool.Pool
}

func NewSettingRepo(db *pgxpool.Pool) *SettingRepo {
	return &SettingRepo{
		db: db,
	}
}

func (r *SettingRepo) CreateUserSetting(
	ctx context.Context,
	user_id int64,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO user_settings (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
		`, user_id,
	)

	return err
}

func (r *SettingRepo) FindUserSettings(
	ctx context.Context,
	user_id int64,
) (models.Settings, error) {
	var settings models.Settings

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
		id,
		user_id,
		memory_reminder,
		location,
		in_app_notifications,
		created_at,
		updated_at
		FROM user_settings
		WHERE user_id = $1
		`, user_id,
	).Scan(
		&settings.Id,
		&settings.User_id,
		&settings.Memory_reminder,
		&settings.Location,
		&settings.In_app_notifications,
		&settings.Created_at,
		&settings.Updated_at,
	)

	if err != nil {
		return models.Settings{}, err
	}

	return settings, nil
}

func (r *SettingRepo) UpdateUserSettings(
	ctx context.Context,
	user_id int64,
	req models.UpdateSetting,
) (models.Settings, error) {
	var settings models.Settings

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE user_settings
		SET
		memory_reminder = COALESCE($2, memory_reminder),
		location = COALESCE($3, location),
		in_app_notifications = COALESCE($4, in_app_notifications),
		updated_at = NOW()
		WHERE user_id = $1
		RETURNING
		id,
		user_id,
		memory_reminder,
		location,
		in_app_notifications,
		created_at,
		updated_at
		`,
		user_id,
		req.Memory_reminder,
		req.Location,
		req.In_app_notifications,
	).Scan(
		&settings.Id,
		&settings.User_id,
		&settings.Memory_reminder,
		&settings.Location,
		&settings.In_app_notifications,
		&settings.Created_at,
		&settings.Updated_at,
	)

	if err != nil {
		return models.Settings{}, err
	}

	return settings, nil
}
