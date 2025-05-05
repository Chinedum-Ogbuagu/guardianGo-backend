package sms

import "gorm.io/gorm"

type Repository interface {
	SaveSMSLog(db *gorm.DB, log *SMSLog) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) SaveSMSLog(db *gorm.DB, log *SMSLog) error {
	return db.Create(log).Error
}
