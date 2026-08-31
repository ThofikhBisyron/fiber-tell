package routers

import (
	"os"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RouterCombine(
	app *fiber.App,
	db *pgxpool.Pool,
) {
	jwtService := services.NewJwtService(
		os.Getenv("JWT_ACCESS_SECRET"),
		os.Getenv("JWT_REFRESH_SECRET"),
	)

	AuthRouters(
		app.Group("/auth"),
		db,
		jwtService,
	)
}
