package guardian

import "gorm.io/gorm"

type Repository interface {
	FindByPhone(db *gorm.DB, phone string) (*Guardian, error)
	Create(db *gorm.DB, g *Guardian) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindByPhone(db *gorm.DB, phone string) (*Guardian, error) {
	var guardian Guardian
	if err := db.Where("phone = ?", phone).First(&guardian).Error; err != nil {
		return nil, err
	}
	return &guardian, nil
}

func (r *repository) Create(db *gorm.DB, g *Guardian) error {
	return db.Create(g).Error
}
