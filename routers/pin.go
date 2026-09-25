package routers

import (
	"tell-be/controllers"
	"tell-be/middlewares"
	"tell-be/repositories"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PinRouters(
	rg fiber.Router,
	db *pgxpool.Pool,
	jwtService *services.JwtService,
) {
	pinRepository := repositories.NewPinRepo(db)

	pinController := controllers.NewPinController(
		pinRepository,
	)

	middleware := rg.Group(
		"/",
		middlewares.AuthMiddleware(jwtService),
	)

	// rg.Post(
	// 	"/pin",
	// 	pinController.CreatePin,
	// )

	middleware.Post("/pin", pinController.CreatePin)
	middleware.Patch("/pin", pinController.UpdatePin)
	middleware.Delete("/pin", pinController.DeletePin)
	middleware.Post("/pin/verify", pinController.VerifyPin)
	middleware.Get("/pin/status", pinController.GetPinStatus)

}
