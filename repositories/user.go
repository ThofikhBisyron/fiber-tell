package repositories

import (
	"context"
	"tell-be/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) FindUserByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	var user models.User

	err := r.db.QueryRow(
		ctx,
		`SELECT
		id,
		email,
		created_at,
		updated_at
		FROM users
		WHERE email = $1`,
		email,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO users (email)
		VALUES ($1)
		RETURNING
		id,
		email,
		created_at,
		updated_at`,
		email,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
