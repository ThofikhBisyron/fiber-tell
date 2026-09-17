package routers

import (
	"tell-be/controllers"
	"tell-be/middlewares"
	"tell-be/repositories"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DiaryRouters(
	rg fiber.Router,
	db *pgxpool.Pool,
	jwtService *services.JwtService,
) {
	diaryRepository := repositories.NewDiaryRepo(db)

	diaryController := controllers.NewDiaryController(
		diaryRepository,
	)

	middleware := rg.Group(
		"/",
		middlewares.AuthMiddleware(jwtService),
	)

	middleware.Get(
		"/recent",
		diaryController.GetRecentDiary,
	)

	middleware.Get(
		"/month",
		diaryController.GetDiaryByMonth,
	)

	middleware.Get(
		"/date/:date",
		diaryController.GetDiaryByDate,
	)

	middleware.Get(
		"/:id",
		diaryController.GetDiaryById,
	)

	middleware.Post(
		"/",
		diaryController.CreateDiary,
	)

	middleware.Patch(
		"/:id",
		diaryController.UpdateDiary,
	)

	middleware.Delete(
		"/:id",
		diaryController.DeleteDiary,
	)

}
