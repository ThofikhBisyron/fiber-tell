package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RouterCombine(
	app *fiber.App,
	db *pgxpool.Pool,
) {
	AuthRouters(
		app.Group("/auth"),
		db,
	)
}
