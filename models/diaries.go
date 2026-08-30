package models

import "time"

type Diaries struct {
	Id         int64     `json:"id"`
	User_id    int64     `json:"user_id" form:"user_id"`
	Mood_id    int64     `json:"mood_id" form:"mood_id"`
	Title      string    `json:"title" form:"title"`
	Content    string    `json:"content" form:"content"`
	Date       time.Time `json:"date" form:"date"`
	Created_at time.Time `json:"created_at" form:"created_at"`
	Updated_at time.Time `json:"updated_at" form:"updated_at"`
}
