package sms

import "time"

type SMSLog struct {
	ID          uint      `gorm:"primaryKey"`
	PhoneNumber string    `gorm:"not null"`
	Message     string    `gorm:"type:text;not null"`
	SentAt      time.Time `gorm:"autoCreateTime"`
}
