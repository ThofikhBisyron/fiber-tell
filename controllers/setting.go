package controllers

import (
	"errors"
	"tell-be/models"
	"tell-be/repositories"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type SettingController struct {
	settingRepo *repositories.SettingRepo
}

func NewSettingController(
	settingRepo *repositories.SettingRepo,
) *SettingController {
	return &SettingController{
		settingRepo: settingRepo,
	}
}

func (c *SettingController) GetUserSettings(
	ctx fiber.Ctx,
) error {
	user_id, ok := ctx.Locals("user_id").(int64)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	settings, err := c.settingRepo.FindUserSettings(
		ctx.Context(),
		user_id,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User settings not found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user settings",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Settings retrieved successfully",
		"data":    settings,
	})
}

func (c *SettingController) UpdateUserSettings(
	ctx fiber.Ctx,
) error {
	user_id, ok := ctx.Locals("user_id").(int64)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var req models.UpdateSetting

	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.Memory_reminder == nil &&
		req.Location == nil &&
		req.In_app_notifications == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "At least one setting must be provided",
		})
	}

	settings, err := c.settingRepo.UpdateUserSettings(
		ctx.Context(),
		user_id,
		req,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User settings not found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user settings",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Settings updated successfully",
		"data":    settings,
	})
}
