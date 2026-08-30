package models

import "time"

type User struct {
	Id         int64     `json:"id"`
	Email      string    `json:"email" form:"email"`
	Created_at time.Time `json:"created_at" form:"created_at"`
	Updated_at time.Time `json:"updated_at" form:"updated_at"`
}
