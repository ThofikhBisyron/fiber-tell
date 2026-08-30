package models

import "time"

type Emailotp struct {
	Id         int64     `json:"id"`
	Email      string    `json:"email" form:"email"`
	Code_hash  string    `json:"code_hash" form:"code_hash"`
	Attempts   int       `json:"attempts" form:"attempts"`
	Created_at time.Time `json:"created_at" form:"created_at"`
	Expires_at time.Time `json:"expires_at" form:"expires_at"`
}
