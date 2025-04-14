package auth

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required"`
}