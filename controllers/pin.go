package controllers

import (
	"tell-be/lib"
	"tell-be/models"
	"tell-be/repositories"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type PinController struct {
	pinRepo *repositories.PinRepo
}

func NewPinController(
	pinRepo *repositories.PinRepo,
) *PinController {
	return &PinController{
		pinRepo: pinRepo,
	}
}

func (c *PinController) CreatePin(
	ctx fiber.Ctx,
) error {
	var input models.CreatePin

	err := ctx.Bind().Body(&input)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if !isValidPin(input.Pin) {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "PIN must be 6 digits",
		})
	}

	userId := ctx.Locals("user_id").(int64)

	_, err = c.pinRepo.FindByUserId(
		ctx.Context(),
		userId,
	)

	if err == nil {
		return ctx.Status(409).JSON(fiber.Map{
			"message": "PIN already exists",
		})
	}

	if err != pgx.ErrNoRows {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to check PIN",
		})
	}

	pinHash, err := lib.HashPassword(input.Pin)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to hash PIN",
		})
	}

	err = c.pinRepo.CreatePin(
		ctx.Context(),
		userId,
		pinHash,
	)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to create PIN",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "PIN created successfully",
	})
}

func (c *PinController) UpdatePin(
	ctx fiber.Ctx,
) error {
	var input models.CreatePin

	err := ctx.Bind().Body(&input)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if !isValidPin(input.Pin) {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "PIN must be 6 digits",
		})
	}

	userId := ctx.Locals("user_id").(int64)

	pinHash, err := lib.HashPassword(input.Pin)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to hash PIN",
		})
	}

	err = c.pinRepo.UpdatePin(
		ctx.Context(),
		userId,
		pinHash,
	)

	if err == pgx.ErrNoRows {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "PIN not found",
		})
	}

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to update PIN",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "PIN updated successfully",
	})
}

func (c *PinController) VerifyPin(
	ctx fiber.Ctx,
) error {
	var input models.VerifyPin

	err := ctx.Bind().Body(&input)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if !isValidPin(input.Pin) {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "PIN must be 6 digits",
		})
	}

	userId := ctx.Locals("user_id").(int64)

	pinHash, err := c.pinRepo.FindByUserId(
		ctx.Context(),
		userId,
	)

	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "PIN not found",
		})
	}

	if !lib.VerifyPassword(input.Pin, pinHash) {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Invalid PIN",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "PIN verified successfully",
	})
}

func (c *PinController) DeletePin(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id").(int64)

	err := c.pinRepo.Delete(
		ctx.Context(),
		userId,
	)

	if err == pgx.ErrNoRows {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "PIN not found",
		})
	}

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to delete PIN",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "PIN deleted successfully",
	})
}

func (c *PinController) GetPinStatus(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	userPin, err := c.pinRepo.FindByUserId(
		ctx.Context(),
		userId.(int64),
	)

	if err != nil || userPin == "" {
		return ctx.JSON(fiber.Map{
			"message": "User pin not found",
			"result":  false,
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User pin found",
		"result":  true,
	})

}

func isValidPin(pin string) bool {
	if len(pin) != 6 {
		return false
	}

	for _, char := range pin {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
