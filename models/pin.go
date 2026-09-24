package models

import "time"

type Pin struct {
	Id         int64     `json:"id"`
	User_id    int64     `json:"user_id" form:"user_id"`
	Pin_hash   string    `json:"pin_hash" form:"pin_hash"`
	Created_at time.Time `json:"created_at" form:"created_at"`
	Updated_at time.Time `json:"updated_at" form:"updated_at"`
}

type CreatePin struct {
	Pin_hash string `json:"pin_hash" form:"pin_hash"`
}

type VerifyPin struct {
	Pin_hash string `json:"pin_hash" form:"pin_hash"`
}
