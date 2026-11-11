package models

import "time"

type Settings struct {
	Id                   int64     `json:"id"`
	User_id              int64     `json:"user_id" form:"user_id"`
	Memory_reminder      bool      `json:"memory_reminder" form:"memory_reminder"`
	Location             bool      `json:"location" form:"location"`
	In_app_notifications bool      `json:"in_app_notifications" form:"in_app_notifications"`
	Created_at           time.Time `json:"created_at" form:"created_at"`
	Updated_at           time.Time `json:"updated_at" form:"updated_at"`
}
