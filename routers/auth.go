package routers

import (
	"os"
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

	profileRepository := repositories.NewProfileRepo(db)

	authRepository := repositories.NewAuthRepository(db)

	emailOtpRepository := repositories.NewOTPRepo(db)

	emailService := services.NewEmailService()

	googleService := services.NewGoogleService(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)

	authService := services.NewAuthService(
		userRepository,
		authRepository,
		profileRepository,
		emailOtpRepository,
		emailService,
		jwtService,
	)

	authController := controllers.NewAuthController(
		authService,
		googleService,
	)

	rg.Post(
		"/email/request", authController.RequestEmailOtp,
	)

	rg.Post(
		"/email/verify", authController.VerifyEmailOtp,
	)

	rg.Post(
		"/refresh",
		authController.RefreshTokenUser,
	)

	rg.Post(
		"/logout",
		authController.Logout,
	)

	rg.Get("/google", authController.GoogleLogin)
	rg.Get("/google/callback", authController.GoogleCallback)
}
