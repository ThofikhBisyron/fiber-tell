package routers

import (
	"os"
	"tell-be/controllers"
	"tell-be/middlewares"
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
	googleConnectService := services.NewGoogleService(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_CONNECT_REDIRECT_URL"),
	)

	facebookService := services.NewFacebookService(
		os.Getenv("FACEBOOK_CLIENT_ID"),
		os.Getenv("FACEBOOK_CLIENT_SECRET"),
		os.Getenv("FACEBOOK_REDIRECT_URL"),
	)

	facebookConnectService := services.NewFacebookService(
		os.Getenv("FACEBOOK_CLIENT_ID"),
		os.Getenv("FACEBOOK_CLIENT_SECRET"),
		os.Getenv("FACEBOOK_CONNECT_REDIRECT_URL"),
	)

	authService := services.NewAuthService(
		userRepository,
		authRepository,
		profileRepository,
		emailOtpRepository,
		facebookService,
		emailService,
		jwtService,
	)

	authController := controllers.NewAuthController(
		authService,
		googleService,
		googleConnectService,
		facebookService,
		facebookConnectService,
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
	rg.Get("/facebook", authController.FacebookLogin)
	rg.Get("/facebook/callback", authController.FacebookCallback)
	rg.Get("/google/connect", middlewares.AuthMiddleware(jwtService), authController.GoogleConnect)
	rg.Get("/google/connect/callback", middlewares.AuthMiddleware(jwtService), authController.GoogleConnectCallback)
	rg.Get("/facebook/connect", middlewares.AuthMiddleware(jwtService), authController.FacebookConnect)
	rg.Get("/facebook/connect/callback", middlewares.AuthMiddleware(jwtService), authController.FacebookConnectCallback)
	rg.Delete("/google/connect", middlewares.AuthMiddleware(jwtService), authController.DisconnectGoogle)
	rg.Delete("/facebook/connect", middlewares.AuthMiddleware(jwtService), authController.DisconnectFacebook)

}
