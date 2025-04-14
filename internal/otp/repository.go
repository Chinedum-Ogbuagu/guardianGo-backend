package otp

import (
	"gorm.io/gorm"
)

type Repository interface {
    SaveOTP(db *gorm.DB, otpRequest *OTPRequest) error
    FindPinIDByPhoneNumber(db *gorm.DB, phoneNumber string) (*OTPRequest, error)
}

type repository struct{}

func NewRepository() Repository {
    return &repository{}
}

func (r *repository) SaveOTP(db *gorm.DB, otpRequest *OTPRequest) error {
    
    return db.Save(otpRequest).Error
}

func (r *repository) FindPinIDByPhoneNumber(db *gorm.DB, phoneNumber string) (*OTPRequest, error) {
    var otpRequest OTPRequest
    if err := db.Where("phone_number = ?", phoneNumber).First(&otpRequest).Error; err != nil {
        return nil, err
    }
    return &otpRequest, nil
}