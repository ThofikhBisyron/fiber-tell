package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PinRepo struct {
	db *pgxpool.Pool
}

func NewPinRepo(db *pgxpool.Pool) *PinRepo {
	return &PinRepo{
		db: db,
	}
}

func (r *PinRepo) FindByUserId(
	ctx context.Context,
	userId int64,
) (string, error) {
	var pinhash string

	err := r.db.QueryRow(
		ctx,
		`
		SELECT pin_hash
		FROM user_security
		WHERE user_id =$1
		`, userId,
	).Scan(&pinhash)

	if err != nil {
		return "", err
	}

	return pinhash, nil
}

func (r *PinRepo) CreatePin(
	ctx context.Context,
	userId int64,
	pinHash string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO user_security 
		(
		user_id,
		pin_hash	
		)
		VALUES ($1, $2)
		`, userId, pinHash,
	)

	return err
}

func (r *PinRepo) UpdatePin(
	ctx context.Context,
	userId int64,
	pinHash string,
) error {
	result, err := r.db.Exec(
		ctx,
		`
		UPDATE user_security
		SET
		pin_hash = $1,
		updated_at = NOW()
		WHERE user_id = $2
		`, pinHash, userId,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *PinRepo) Delete(
	ctx context.Context,
	userId int64,
) error {
	result, err := r.db.Exec(
		ctx,
		`
		DELETE FROM user_security
		WHERE user_id = $1
		`,
		userId,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
