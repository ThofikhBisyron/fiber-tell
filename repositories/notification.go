package repositories

import (
	"context"
	"tell-be/models"

	"github.com/jackc/pgx/v5"
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
	notification models.Notification,
) (models.Notification, error) {

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO notifications (
			user_id,
			type,
			title,
			message,
			reference_id
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			user_id,
			type,
			title,
			message,
			reference_id,
			is_read,
			created_at
		`,
		notification.User_id,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Reference_id,
	).Scan(
		&notification.Id,
		&notification.User_id,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.Reference_id,
		&notification.Is_read,
		&notification.Created_at,
	)

	if err != nil {
		return models.Notification{}, err
	}

	return notification, nil
}

func (r *NotifRepo) FindNotifications(
	ctx context.Context,
	user_id int64,
	limit int,
	offset int,
) ([]models.Notification, error) {

	rows, err := r.db.Query(
		ctx,
		`
        SELECT
            id,
            user_id,
            type,
            title,
            message,
            reference_id,
            is_read,
            created_at
        FROM notifications
        WHERE user_id = $1
        ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
        `,
		user_id,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notifications := make([]models.Notification, 0)

	for rows.Next() {
		var notification models.Notification

		err := rows.Scan(
			&notification.Id,
			&notification.User_id,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.Reference_id,
			&notification.Is_read,
			&notification.Created_at,
		)

		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotifRepo) MarkAsRead(
	ctx context.Context,
	user_id int64,
	notification_id int64,
) error {

	result, err := r.db.Exec(
		ctx,
		`
        UPDATE notifications
        SET is_read = TRUE
        WHERE id = $1
        AND user_id = $2
        `,
		notification_id,
		user_id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *NotifRepo) MarkAllAsRead(
	ctx context.Context,
	user_id int64,
) error {

	_, err := r.db.Exec(
		ctx,
		`
        UPDATE notifications
        SET is_read = TRUE
        WHERE user_id = $1
        AND is_read = FALSE
        `,
		user_id,
	)

	return err
}

func (r *NotifRepo) DeleteNotification(
	ctx context.Context,
	user_id int64,
	notification_id int64,
) error {

	result, err := r.db.Exec(
		ctx,
		`
        DELETE FROM notifications
        WHERE id = $1
        AND user_id = $2
        `,
		notification_id,
		user_id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *NotifRepo) CountUnread(
	ctx context.Context,
	user_id int64,
) (int64, error) {

	var count int64

	err := r.db.QueryRow(
		ctx,
		`
        SELECT COUNT(*)
        FROM notifications
        WHERE user_id = $1
        AND is_read = FALSE
        `,
		user_id,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
