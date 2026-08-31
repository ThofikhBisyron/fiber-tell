package middlewares

import (
	"strconv"

	"tell-be/services"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(jwtService *services.JwtService) fiber.Handler {
	return func(c fiber.Ctx) error {

		tokenString := c.Cookies("access_token")

		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid or expired token",
			})
		}

		userIDString, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid user id",
			})
		}

		userID, err := strconv.ParseInt(userIDString, 10, 64)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid user id",
			})
		}

		c.Locals("user_id", userID)

		return c.Next()
	}
}
