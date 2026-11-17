package routers

import (
	"tell-be/controllers"
	"tell-be/middlewares"
	"tell-be/repositories"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SettingRouters(
	rg fiber.Router,
	db *pgxpool.Pool,
	jwtService *services.JwtService,
) {
	settingRepo := repositories.NewSettingRepo(db)

	settingController := controllers.NewSettingController(
		settingRepo,
	)

	middleware := rg.Group("/", middlewares.AuthMiddleware(jwtService))
	middleware.Get("/", settingController.GetUserSettings)
	middleware.Patch("/", settingController.UpdateUserSettings)

	// rg.Get("/",middlewares.AuthMiddleware(jwtService), settingController.GetUserSettings,)

	// rg.Patch("/",middlewares.AuthMiddleware(jwtService), settingController.UpdateUserSettings,)
}
