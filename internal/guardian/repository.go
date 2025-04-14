package guardian

import "gorm.io/gorm"

type Repository interface {
	Create(db *gorm.DB, guardian *Guardian ) error 
	GetByPhoneNumber(db *gorm.DB, phone string) (*Guardian, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, guardian *Guardian) error {
	return db.Create(guardian).Error
}

func (r *repository) GetByPhoneNumber(db *gorm.DB, phone string ) (*Guardian, error) {
	var p Guardian 
	if err := db.Where("phone_number = ?", phone).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

