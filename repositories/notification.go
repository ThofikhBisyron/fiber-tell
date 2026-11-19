package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotifRepo struct {
	db *pgxpool.Pool
}

func NewNotifRepo(db *pgxpool.Pool) *NotifRepo {
	return &NotifRepo{
		db: db,
	}
}

func (r *NotifRepo) CreateNotification(
	ctx context.Context,
	user_id int64,
) error {

}
