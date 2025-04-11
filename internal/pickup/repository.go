package pickup

import "gorm.io/gorm"

type Repository interface {
	Create(db *gorm.DB, p *Pickup) error
	GetByID(db *gorm.DB, id uint) (*Pickup, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, p *Pickup) error {
	return db.Create(p).Error
}

func (r *repository) GetByID(db *gorm.DB, id uint) (*Pickup, error) {
	var p Pickup
	if err := db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}