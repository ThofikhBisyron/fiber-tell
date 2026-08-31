package routers

import (
	"tell-be/controllers"
	"tell-be/repositories"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRouters(
	rg fiber.Router,
	db *pgxpool.Pool,
	jwtService *services.JwtService,
) {

	userRepository := repositories.NewUserRepository(db)

	emailOtpRepository := repositories.NewOTPRepo(db)

	emailService := services.NewEmailService()

	authService := services.NewAuthService(
		userRepository,
		emailOtpRepository,
		emailService,
		jwtService,
	)

	authController := controllers.NewAuthController(
		authService,
	)

	rg.Post(
		"/email/request", authController.RequestEmailOtp,
	)

	rg.Post(
		"/email/verify", authController.VerifyEmailOtp,
	)
}
