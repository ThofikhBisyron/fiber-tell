package models

import "time"

type Auth struct {
	Id               int64     `json:"id"`
	User_id          int64     `json:"user_id" form:"user_id"`
	Provider_id      int64     `json:"provider_id" form:"provider_id"`
	Provider_user_id string    `json:"provider_user_id" form:"provider_user_id"`
	Created_at       time.Time `json:"created_at" form:"created_at"`
	Updated_at       time.Time `json:"updated_at" form:"updated_at"`
}
