package controllers

import "tell-be/repositories"

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
