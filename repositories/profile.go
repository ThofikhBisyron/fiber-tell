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
) (models.JoinEmailAuthProfile, error) {
	var profile models.JoinEmailAuthProfile

	err := r.db.QueryRow(
		ctx,
		`SELECT
			u.email AS email_user,
			COALESCE(
				MAX(ua.provider_email) FILTER (
					WHERE up.name = 'google'
				), ''
			) AS email_google,
			COALESCE(
				MAX(ua.provider_email) FILTER (
					WHERE up.name = 'facebook'
				), ''
			) AS email_facebook,
			p.first_name,
			p.last_name,
			p.phone_number
		FROM users u
		INNER JOIN user_auth ua
			ON u.id = ua.user_id
		INNER JOIN user_provider up
			ON ua.provider_id = up.id
		INNER JOIN profiles p
			ON u.id = p.user_id
		WHERE u.id = $1
		GROUP BY u.id, u.email, p.first_name, p.last_name, p.phone_number;`,
		user_id,
	).Scan(
		&profile.Email_user,
		&profile.Email_google,
		&profile.Email_facebook,
		&profile.First_name,
		&profile.Last_name,
		&profile.Phone_number,
	)

	if err != nil {
		return models.JoinEmailAuthProfile{}, err
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
