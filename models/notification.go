package models

import "time"

type Notification struct {
	Id           int64     `json:"id"`
	User_id      int64     `json:"user_id" form:"user_id"`
	Type         string    `json:"type" form:"type"`
	Title        string    `json:"title" form:"title"`
	Message      string    `json:"message" form:"message"`
	Reference_id *int64    `json:"reference_id" form:"reference_id"`
	Is_read      bool      `json:"is_read" form:"is_read"`
	Created_at   time.Time `json:"created_at" form:"created_at"`
}
