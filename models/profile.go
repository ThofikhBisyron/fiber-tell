package models

import "time"

type Profile struct {
	Id           int64     `json:"id"`
	User_id      int64     `json:"user_id" form:"user_id"`
	First_name   string    `json:"first_name" form:"first_name"`
	Last_name    string    `json:"last_name" form:"last_name"`
	Phone_number string    `json:"phone_number" form:"phone_number"`
	Created_at   time.Time `json:"created_at" form:"created_at"`
	Updated_at   time.Time `json:"updated_at" form:"updated_at"`
}
