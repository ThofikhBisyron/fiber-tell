package controllers

import (
	"strconv"
	"tell-be/repositories"
	"time"

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

func (c *StatController) GetStatistics(
	ctx fiber.Ctx,
) error {
	userId, ok := ctx.Locals("user_id").(int64)

	if !ok {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	period := ctx.Query("period")
	dateQuery := ctx.Query("date")

	if period == "" {
		period = "weekly"
	}

	if dateQuery == "" {
		dateQuery = time.Now().UTC().Format("2006-01-02")
	}

	selectedDate, err := time.Parse(
		"2006-01-02",
		dateQuery,
	)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	var startDate time.Time
	var endDate time.Time

	switch period {
	case "weekly":
		weekday := int(selectedDate.Weekday())

		daysSinceMonday := (weekday + 6) % 7

		startDate = selectedDate.AddDate(
			0,
			0,
			-daysSinceMonday,
		)

		endDate = startDate.AddDate(0, 0, 7)

	case "monthly":
		startDate = time.Date(
			selectedDate.Year(),
			selectedDate.Month(),
			1,
			0, 0, 0, 0,
			selectedDate.Location(),
		)

		endDate = startDate.AddDate(0, 1, 0)

	case "yearly":
		startDate = time.Date(
			selectedDate.Year(),
			time.January,
			1,
			0, 0, 0, 0,
			selectedDate.Location(),
		)

		endDate = startDate.AddDate(1, 0, 0)

	default:
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid period. Use weekly, monthly, or yearly",
		})
	}

	totalDiary, err := c.statRepo.CountDiaryByPeriod(
		ctx.Context(),
		userId,
		startDate,
		endDate,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to get total diary statistics",
		})
	}

	moodCounts, err := c.statRepo.CountDiaryByMoodAndPeriod(
		ctx.Context(),
		userId,
		startDate,
		endDate,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to get mood statistics",
		})
	}

	activity, err := c.statRepo.CountDiaryActivity(
		ctx.Context(),
		userId,
		startDate,
		endDate,
		period,
	)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to get diary activity",
		})
	}

	return ctx.JSON(fiber.Map{
		"period":      period,
		"start_date":  startDate.Format("2006-01-02"),
		"end_date":    endDate.AddDate(0, 0, -1).Format("2006-01-02"),
		"total_diary": totalDiary,
		"mood_counts": moodCounts,
		"activity":    activity,
	})
}
