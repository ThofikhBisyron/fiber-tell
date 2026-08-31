package controllers

import (
	"strings"
	"tell-be/dtos"
	"tell-be/services"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(
	authService *services.AuthService,
) *AuthController {

	return &AuthController{
		authService: authService,
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
