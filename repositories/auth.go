package repositories

import (
	"context"
	"tell-be/models"

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
	provider_email string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO user_auth (
		user_id,
		provider_id,
		provider_user_id,
		provider_email
		)
		VALUES ($1, $2, $3, $4)
		`, user_id, provider_id, provider_user_id, provider_email,
	)

	return err
}

func (r *AuthRepo) FindUserAuthAndProvider(
	ctx context.Context,
	user_id int64,
	provider_id int64,
) (models.AuthUserProvider, error) {
	var auth models.AuthUserProvider

	err := r.db.QueryRow(
		ctx,
		`
		SELECT 
		id,
		user_id,
		provider_id
		FROM user_auth
		WHERE user_id = $1
		AND provider_id = $2
		`, user_id, provider_id,
	).Scan(
		&auth.Id,
		&auth.User_id,
		&auth.Provider_id,
	)

	if err != nil {
		return models.AuthUserProvider{}, err
	}

	return auth, err
}

func (r *AuthRepo) DeleteUserAuth(
	ctx context.Context,
	user_id int64,
	provider_id int64,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM user_auth
		WHERE user_id = $1
		AND provider_id = $2
		`, user_id, provider_id,
	)

	return err
}
