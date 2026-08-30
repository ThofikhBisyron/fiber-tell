package dtos

type RequestEmailOTP struct {
	Email string `json:"email" form:"email"`
}

type VerifyEmailOTP struct {
	Email string `json:"email" form:"email"`
	Code  string `json:"code" form:"code"`
}
