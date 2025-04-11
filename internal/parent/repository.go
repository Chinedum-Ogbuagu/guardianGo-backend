package parent

import "gorm.io/gorm"

type Repository interface {
	Create(db *gorm.DB, parent *Parent ) error 
	GetByPhoneNumber(db *gorm.DB, phone string) (*Parent, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, parent *Parent) error {
	return db.Create(parent).Error
}

func (r *repository) GetByPhoneNumber(db *gorm.DB, phone string ) (*Parent, error) {
	var p Parent 
	if err := db.Where("phone_number = ?", phone).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

