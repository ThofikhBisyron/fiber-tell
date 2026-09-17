package controllers

import (
	"strconv"
	"tell-be/models"
	"tell-be/repositories"
	"time"

	"github.com/gofiber/fiber/v3"
)

type DiaryController struct {
	diaryRepo *repositories.DiaryRepo
}

func NewDiaryController(
	diaryRepo *repositories.DiaryRepo,
) *DiaryController {
	return &DiaryController{
		diaryRepo: diaryRepo,
	}
}

func (c *DiaryController) CreateDiary(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var input models.CreateDiary

	err := ctx.Bind().Body(&input)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid Request Body",
		})
	}

	diary, err := c.diaryRepo.CreateDiary(
		ctx.Context(),
		userId.(int64),
		input.Mood_id,
		input.Title,
		input.Content,
		input.Date,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to create diary",
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "Diary created successfully",
		"diary":   diary,
	})
}

func (c *DiaryController) GetRecentDiary(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	diary, err := c.diaryRepo.FindRecentDiary(
		ctx.Context(),
		userId.(int64),
	)

	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Diary not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"diary": diary,
	})
}

func (c *DiaryController) GetDiaryByMonth(
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
			"message": "Invalid Year",
		})
	}

	month, err := strconv.Atoi(ctx.Query("month"))
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid Month",
		})
	}

	if month < 1 || month > 12 {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Month must be between 1 and 12",
		})
	}

	diaries, err := c.diaryRepo.FindDiariesByMonth(
		ctx.Context(),
		userId.(int64),
		year,
		month,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to get diaries",
		})
	}

	return ctx.JSON(fiber.Map{
		"diaries": diaries,
	})
}

func (c *DiaryController) GetDiaryByDate(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	dateString := ctx.Params("date")

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	diaries, err := c.diaryRepo.FindDiariesByDate(
		ctx.Context(),
		userId.(int64),
		date,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to get diaries",
		})
	}

	return ctx.JSON(fiber.Map{
		"diaries": diaries,
	})
}

func (c *DiaryController) GetDiaryById(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	diaryId, err := strconv.ParseInt(
		ctx.Params("id"),
		10,
		64,
	)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid diary id",
		})
	}

	diary, err := c.diaryRepo.FindDiaryById(
		ctx.Context(),
		userId.(int64),
		diaryId,
	)

	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Diary not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"diary": diary,
	})
}

func (c *DiaryController) UpdateDiary(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	diaryId, err := strconv.ParseInt(
		ctx.Params("id"),
		10,
		64,
	)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid diary id",
		})
	}

	var input models.UpdateDiary

	err = ctx.Bind().Body(&input)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err = c.diaryRepo.UpdateDiaryById(
		ctx.Context(),
		diaryId,
		userId.(int64),
		input.Mood_id,
		input.Title,
		input.Content,
		input.Date,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to update diary",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Diary updated successfully",
	})
}

func (c *DiaryController) DeleteDiary(
	ctx fiber.Ctx,
) error {
	userId := ctx.Locals("user_id")

	if userId == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	diaryId, err := strconv.ParseInt(
		ctx.Params("id"),
		10,
		64,
	)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid diary id",
		})
	}

	err = c.diaryRepo.DeleteDiaryById(
		ctx.Context(),
		diaryId,
		userId.(int64),
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to delete diary",
		})
	}
	return ctx.JSON(fiber.Map{
		"message": "Diary deleted successfully",
	})
}
