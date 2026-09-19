package routers

import (
	"tell-be/controllers"
	"tell-be/middlewares"
	"tell-be/repositories"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StatisticsRouters(
	rg fiber.Router,
	db *pgxpool.Pool,
	jwtService *services.JwtService,
) {
	statRepo := repositories.NewStatRepo(db)

	statController := controllers.NewStatController(
		statRepo,
	)

	middleware := rg.Group(
		"/",
		middlewares.AuthMiddleware(jwtService),
	)

	middleware.Get(
		"/total",
		statController.GetTotalDiary,
	)

	middleware.Get(
		"/mood",
		statController.GetDiaryByMood,
	)

	middleware.Get(
		"/monthly",
		statController.GetDiaryByMonth,
	)

}
