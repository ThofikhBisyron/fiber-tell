package repositories

import (
	"context"

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
