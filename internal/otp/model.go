package otp

type OTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	PinID       string `json:"pin_id"`
	Purpose     string `json:"purpose"`
}