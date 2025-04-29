package otp

import "gorm.io/gorm"

type Repository interface {
	SaveOTP(db *gorm.DB, otp *OTPRequest) error
	FindPinIDByPhoneNumber(db *gorm.DB, phone string, purpose string) (*OTPRequest, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) SaveOTP(db *gorm.DB, otp *OTPRequest) error {
	return db.Save(otp).Error
}

func (r *repository) FindPinIDByPhoneNumber(db *gorm.DB, phone string, purpose string) (*OTPRequest, error) {
	var req OTPRequest
	err := db.Where("phone_number = ?", phone).Order("created_at DESC").First(&req).Error
	return &req, err
}
