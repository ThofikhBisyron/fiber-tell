package controllers

import (
	"errors"
	"strconv"
	"tell-be/repositories"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type NotificationController struct {
	notifRepo *repositories.NotifRepo
}

func NewNotificationController(
	notifRepo *repositories.NotifRepo,
) *NotificationController {
	return &NotificationController{
		notifRepo: notifRepo,
	}
}

// GET /notifications
func (c *NotificationController) GetNotifications(
	ctx fiber.Ctx,
) error {
	userId, ok := ctx.Locals("user_id").(int64)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"message": "Unauthorized",
			},
		)
	}

	limit := 20
	offset := 0

	if query := ctx.Query("limit"); query != "" {
		value, err := strconv.Atoi(query)

		if err != nil || value < 1 || value > 100 {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"message": "Limit must be between 1 and 100",
				},
			)
		}

		limit = value
	}

	if query := ctx.Query("offset"); query != "" {
		value, err := strconv.Atoi(query)

		if err != nil || value < 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"message": "Invalid offset",
				},
			)
		}

		offset = value
	}

	notifications, err := c.notifRepo.FindNotifications(
		ctx.Context(),
		userId,
		limit,
		offset,
	)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"message": "Failed to get notifications",
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"notifications": notifications,
			"limit":         limit,
			"offset":        offset,
		},
	)
}

// GET /notifications/unread/count
func (c *NotificationController) GetUnreadCount(
	ctx fiber.Ctx,
) error {
	userId, ok := ctx.Locals("user_id").(int64)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"message": "Unauthorized",
			},
		)
	}

	count, err := c.notifRepo.CountUnread(
		ctx.Context(),
		userId,
	)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"message": "Failed to count unread notifications",
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"unread_count": count,
		},
	)
}

// PATCH /notifications/:id/read
func (c *NotificationController) MarkAsRead(
	ctx fiber.Ctx,
) error {
	userId, ok := ctx.Locals("user_id").(int64)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"message": "Unauthorized",
			},
		)
	}

	notificationId, err := strconv.ParseInt(
		ctx.Params("id"),
		10,
		64,
	)

	if err != nil || notificationId <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": "Invalid notification id",
			},
		)
	}

	err = c.notifRepo.MarkAsRead(
		ctx.Context(),
		userId,
		notificationId,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ctx.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"message": "Notification not found",
			},
		)
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"message": "Failed to mark notification as read",
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "Notification marked as read",
		},
	)
}

// PATCH /notifications/read-all
func (c *NotificationController) MarkAllAsRead(
	ctx fiber.Ctx,
) error {
	userId, ok := ctx.Locals("user_id").(int64)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"message": "Unauthorized",
			},
		)
	}

	err := c.notifRepo.MarkAllAsRead(
		ctx.Context(),
		userId,
	)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"message": "Failed to mark all notifications as read",
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "All notifications marked as read",
		},
	)
}

// DELETE /notifications/:id
func (c *NotificationController) DeleteNotification(
	ctx fiber.Ctx,
) error {
	userId, ok := ctx.Locals("user_id").(int64)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"message": "Unauthorized",
			},
		)
	}

	notificationId, err := strconv.ParseInt(
		ctx.Params("id"),
		10,
		64,
	)

	if err != nil || notificationId <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": "Invalid notification id",
			},
		)
	}

	err = c.notifRepo.DeleteNotification(
		ctx.Context(),
		userId,
		notificationId,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ctx.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"message": "Notification not found",
			},
		)
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"message": "Failed to delete notification",
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "Notification deleted successfully",
		},
	)
}
