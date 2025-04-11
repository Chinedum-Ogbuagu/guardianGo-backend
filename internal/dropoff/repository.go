package dropoff

import "gorm.io/gorm"


type Repository interface {
	Create(db *gorm.DB, d *DropOff) error
	GetByID(db *gorm.DB, id uint) (*DropOff, error)
}

type repository struct {}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, d *DropOff) error {
	return db.Create(d).Error
}

func (r *repository) GetByID(db *gorm.DB, id uint) (*DropOff, error) {
	var d DropOff
	if err := db.First(&d, id).Error; err != nil {
		return nil, err
	}

	return &d, nil
}