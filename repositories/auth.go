package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepository(
	db *pgxpool.Pool,
) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) FindProviderByName(
	ctx context.Context,
	name string,
) (int64, error) {
	var provider_id int64

	err := r.db.QueryRow(
		ctx,
		`
		SELECT id
		FROM user_provider
		WHERE name = $1
		`, name,
	).Scan(&provider_id)

	if err != nil {
		return 0, err
	}

	return provider_id, nil
}

func (r *AuthRepo) FindUserByProvider(
	ctx context.Context,
	provider_id int64,
	provider_user_id string,
) (int64, error) {
	var user_id int64

	err := r.db.QueryRow(
		ctx,
		`
		SELECT user_id
		FROM user_auth
		WHERE provider_id = $1
		AND provider_user_id = $2
		`, provider_id, provider_user_id,
	).Scan(&user_id)

	if err != nil {
		return 0, err
	}

	return user_id, nil
}

func (r *AuthRepo) CreateUserAuth(
	ctx context.Context,
	user_id int64,
	provider_id int64,
	provider_user_id string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO user_auth (
		user_id,
		provider_id,
		provider_user_id
		)
		VALUES ($1, $2, $3)
		`, user_id, provider_id, provider_user_id,
	)

	return err
}
