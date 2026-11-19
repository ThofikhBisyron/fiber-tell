package routers

import (
	"tell-be/controllers"
	"tell-be/middlewares"
	"tell-be/repositories"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NotifRouters(
	rg fiber.Router,
	db *pgxpool.Pool,
	jwtService *services.JwtService,
) {
	notifRepository := repositories.NewNotifRepo(db)

	notifController := controllers.NewNotificationController(
		notifRepository,
	)

	middleware := rg.Group(
		"/",
		middlewares.AuthMiddleware(jwtService),
	)

	middleware.Get("/", notifController.GetNotifications)

	middleware.Get(
		"/unread/count",
		notifController.GetUnreadCount,
	)

	middleware.Patch(
		"/read-all",
		notifController.MarkAllAsRead,
	)

	middleware.Patch(
		"/:id/read",
		notifController.MarkAsRead,
	)

	middleware.Delete(
		"/:id",
		notifController.DeleteNotification,
	)
}
