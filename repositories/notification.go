package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type NotifRepo struct {
	db *pgxpool.Pool
}

func NewNotifRepo(db *pgxpool.Pool) *NotifRepo {
	return &NotifRepo{
		db: db,
	}
}
