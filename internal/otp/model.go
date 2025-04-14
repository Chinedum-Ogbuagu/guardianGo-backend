package otp

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" gorm:"primaryKey"`
	PinID       string `json:"pin_id"`
	Purpose     string `json:"purpose"`
	DropOffID   uint   `json:"drop_off_id"`
}