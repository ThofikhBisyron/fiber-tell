package routers

import (
	"tell-be/controllers"
	"tell-be/middlewares"
	"tell-be/repositories"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProfileRouters(
	rg fiber.Router,
	db *pgxpool.Pool,
	jwtService *services.JwtService,
) {
	profileRepository := repositories.NewProfileRepo(db)

	profileController := controllers.NewProfileController(
		profileRepository,
	)

	protected := rg.Group("/", middlewares.AuthMiddleware(jwtService))

	protected.Get(
		"/me",
		profileController.GetProfile,
	)

	protected.Patch(
		"/",
		profileController.UpdateProfile,
	)
	// rg.Get(
	// 	"/",
	// 	middlewares.AuthMiddleware(jwtService),
	// 	profileController.GetProfile,
	// )

	// rg.Put(
	// 	"/",
	// 	middlewares.AuthMiddleware(jwtService),
	// 	profileController.UpdateProfile,
	// )
}
