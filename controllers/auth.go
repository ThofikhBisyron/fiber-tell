package controllers

import (
	"strings"
	"tell-be/dtos"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthController struct {
	authService   *services.AuthService
	googleService *services.GoogleService
}

func NewAuthController(
	authService *services.AuthService,
	googleService *services.GoogleService,
) *AuthController {

	return &AuthController{
		authService:   authService,
		googleService: googleService,
	}
}

func (c *AuthController) RequestEmailOtp(
	ctx fiber.Ctx,
) error {
	var input dtos.RequestEmailOTP

	if err := ctx.Bind().Body(&input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid Request Body",
		})
	}

	input.Email = strings.ToLower(
		strings.TrimSpace(input.Email),
	)

	if input.Email == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Email is required",
		})
	}

	err := c.authService.RequestOtp(
		ctx.Context(),
		input.Email,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Otp sent successfully",
	})
}

func (c *AuthController) VerifyEmailOtp(
	ctx fiber.Ctx,
) error {
	var input dtos.VerifyEmailOTP

	if err := ctx.Bind().Body(&input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	input.Email = strings.ToLower(
		strings.TrimSpace(input.Email),
	)

	if input.Email == "" || input.Code == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Email and OTP required",
		})
	}

	result, err := c.authService.VerifyOtp(
		ctx.Context(),
		input.Email,
		input.Code,
	)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   15 * 60,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
	})

	return ctx.JSON(fiber.Map{
		"message": "Login successful",
		"user":    result.User,
	})
}

func (c *AuthController) RefreshTokenUser(
	ctx fiber.Ctx,
) error {
	refreshToken := ctx.Cookies("refresh_token")

	if refreshToken == "" {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Refresh token not found",
		})
	}

	result, err := c.authService.RefreshToken(
		ctx.Context(),
		refreshToken,
	)

	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   15 * 60,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
	})

	return ctx.JSON(fiber.Map{
		"message": "Token refreshed successfully",
	})

}

func (c *AuthController) Logout(
	ctx fiber.Ctx,
) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   -1,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   -1,
	})

	return ctx.JSON(fiber.Map{
		"message": "Logout succesful",
	})
}

func (c *AuthController) GoogleLogin(
	ctx fiber.Ctx,
) error {
	state := uuid.NewString()

	ctx.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   10 * 60,
	})

	url := c.googleService.GetAuthURL(state)

	return ctx.Redirect().To(url)
}

func (c *AuthController) GoogleCallback(
	ctx fiber.Ctx,
) error {
	state := ctx.Query("state")
	code := ctx.Query("code")

	if state == "" || code == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid Google OAuth callback",
		})
	}

	oauthState := ctx.Cookies("oauth_state")

	if oauthState == "" || oauthState != state {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid OAuth state",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   -1,
	})

	googleUser, err := c.googleService.GetUser(
		ctx.Context(),
		code,
	)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	result, err := c.authService.LoginWithGoogle(
		ctx.Context(),
		googleUser,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   15 * 60,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
	})

	return ctx.JSON(fiber.Map{
		"message": "Google login successful",
		"user":    result.User,
	})
}
