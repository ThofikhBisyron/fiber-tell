package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type SettingRepo struct {
	db *pgxpool.Pool
}

func NewSettingRepo(db *pgxpool.Pool) *SettingRepo {
	return &SettingRepo{
		db: db,
	}
}
