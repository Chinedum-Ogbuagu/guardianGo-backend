package otp

import "time"

type OTPRequest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PhoneNumber string    `json:"phone_number" gorm:"index"` // Make this indexed but not primary key
	PinID       string    `json:"pin_id"`
	Purpose     string    `json:"purpose"`
	DropOffID   uint      `json:"drop_off_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type OTPToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"index" json:"phone"`
	Purpose   string    `json:"purpose"`
	DropOffID uint      `json:"drop_off_id"`
	PinID     string    `json:"pin_id"`     // from Termii
	ExpiresAt time.Time `json:"expires_at"` // Optional, just in case
	CreatedAt time.Time `json:"created_at"`
}
