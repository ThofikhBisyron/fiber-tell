package controllers

import (
	"tell-be/models"
	"tell-be/repositories"

	"github.com/gofiber/fiber/v3"
)

type ProfileController struct {
	profileRepo *repositories.ProfileRepo
}

func NewProfileController(
	profileRepo *repositories.ProfileRepo,
) *ProfileController {
	return &ProfileController{
		profileRepo: profileRepo,
	}
}

func (c *ProfileController) GetProfile(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	profile, err := c.profileRepo.FindProfileByUserId(
		ctx.Context(),
		userId.(int64),
	)

	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Profile Not Found",
		})
	}

	return ctx.JSON(fiber.Map{
		"profile": profile,
	})
}

func (c *ProfileController) UpdateProfile(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var input models.UpdateProfile

	if err := ctx.Bind().Body(&input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"Message": "Invalid request body",
		})
	}

	err := c.profileRepo.UpdateProfileByUserId(
		ctx.Context(),
		input.First_name,
		input.Last_name,
		input.Phone_number,
		userId.(int64),
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Wait, Maintenance",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}
