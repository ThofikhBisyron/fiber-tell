package repositories

import (
	"context"
	"tell-be/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepo struct {
	db *pgxpool.Pool
}

func NewProfileRepo(db *pgxpool.Pool) *ProfileRepo {
	return &ProfileRepo{
		db: db,
	}
}

func (r *ProfileRepo) CreateProfile(
	ctx context.Context,
	user_id int64,
) (models.Profile, error) {
	var profile models.Profile

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO profiles (user_id)
		VALUES ($1)
		RETURNING
		id,
		user_id,
		first_name,
		last_name,
		phone_number,
		created_at,
		updated_at
		`,
		user_id,
	).Scan(
		&profile.Id,
		&profile.User_id,
		&profile.First_name,
		&profile.Last_name,
		&profile.Phone_number,
		&profile.Created_at,
		&profile.Updated_at,
	)

	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (r *ProfileRepo) FindProfileByUserId(
	ctx context.Context,
	user_id int64,
) (models.Profile, error) {
	var profile models.Profile

	err := r.db.QueryRow(
		ctx,
		`SELECT
		id,
		user_id,
		first_name,
		last_name,
		phone_number,
		created_at,
		updated_at
		FROM profiles
		WHERE user_id = ($1)`,
		user_id,
	).Scan(
		&profile.Id,
		&profile.User_id,
		&profile.First_name,
		&profile.Last_name,
		&profile.Phone_number,
		&profile.Created_at,
		&profile.Updated_at,
	)

	if err != nil {
		return models.Profile{}, err
	}

	return profile, err
}

func (r *ProfileRepo) UpdateProfileByUserId(
	ctx context.Context,
	first_name string,
	last_name string,
	phone_number string,
	user_id int64,
) error {

	_, err := r.db.Exec(
		ctx,
		`UPDATE profiles
		SET
		first_name = $1,
		last_name = $2,
		phone_number = $3,
		updated_at = NOW()
		WHERE user_id = $4`, first_name, last_name, phone_number, user_id,
	)

	return err
}
