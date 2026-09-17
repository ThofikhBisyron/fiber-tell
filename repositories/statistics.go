package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type StatRepo struct {
	db *pgxpool.Pool
}

func NewStatRepo(db *pgxpool.Pool) *StatRepo {
	return &StatRepo{
		db: db,
	}
}
