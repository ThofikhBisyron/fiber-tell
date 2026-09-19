package controllers

import (
	"strconv"
	"tell-be/repositories"

	"github.com/gofiber/fiber/v3"
)

type StatController struct {
	statRepo *repositories.StatRepo
}

func NewStatController(
	statRepo *repositories.StatRepo,
) *StatController {
	return &StatController{
		statRepo: statRepo,
	}
}

func (c *StatController) GetTotalDiary(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	total, err := c.statRepo.CountDiary(
		ctx.Context(),
		userId.(int64),
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to get statistics",
		})
	}

	return ctx.JSON(fiber.Map{
		"total_diary": total,
	})
}

func (c *StatController) GetDiaryByMood(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	result, err := c.statRepo.CountDiaryByMood(
		ctx.Context(),
		userId.(int64),
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to get mood statistics",
		})
	}

	return ctx.JSON(fiber.Map{
		"mood_counts": result,
	})
}

func (c *StatController) GetDiaryByMonth(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	year, err := strconv.Atoi(ctx.Query("year"))

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid year",
		})
	}

	result, err := c.statRepo.CountDiaryByMonth(
		ctx.Context(),
		userId.(int64),
		year,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to get monthly statistics",
		})
	}

	return ctx.JSON(fiber.Map{
		"year":         year,
		"monthly_data": result,
	})
}
