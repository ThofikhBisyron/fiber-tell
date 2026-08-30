package repositories

import (
	"context"
	"tell-be/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OTPRepo struct {
	db *pgxpool.Pool
}

func NewOTPRepo(db *pgxpool.Pool) *OTPRepo {
	return &OTPRepo{
		db: db,
	}
}

func (r *OTPRepo) DeleteOtpByEmail(
	ctx context.Context,
	email string,
) error {
	_, err := r.db.Exec(
		ctx,
		`DELETE from email_otps WHERE email = $1`, email,
	)

	return err
}

func (r *OTPRepo) CreateOtp(
	ctx context.Context,
	email string,
	code_hash string,
	expires_at time.Time,
) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO email_otps(
		email,
		code_hash,
		expires_at
		)
		VALUES ($1, $2, $3)`, email, code_hash, expires_at,
	)

	return err
}

func (r *OTPRepo) FindOtpByEmail(
	ctx context.Context,
	email string,
) (models.Emailotp, error) {
	var otp models.Emailotp

	err := r.db.QueryRow(
		ctx,
		`SELECT 
		id,
		email,
		code_hash,
		attempts,
		created_at,
		expires_at
		FROM email_otps
		WHERE email = $1
		ORDER BY created_at DESC
		LIMIT 1`, email,
	).Scan(
		&otp.Id,
		&otp.Email,
		&otp.Code_hash,
		&otp.Attempts,
		&otp.Created_at,
		&otp.Expires_at,
	)

	if err != nil {
		return models.Emailotp{}, err
	}

	return otp, nil
}

func (r *OTPRepo) IncreaseAttempts(
	ctx context.Context,
	id int,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE email_otps
		SET attempts = attempts + 1
		WHERE id = $1
		`, id,
	)

	return err
}

func (r *OTPRepo) DeleteOtpById(
	ctx context.Context,
	id int64,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM email_otps
		WHERE id = $1`, id,
	)

	return err
}
