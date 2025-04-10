package church

import (
	"gorm.io/gorm"
)

type Repository interface {
	CreateChurch(db *gorm.DB, church *Church) error
	GetChurchById(db *gorm.DB, id uint) (*Church, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateChurch(db *gorm.DB, church *Church) error {
	return db.Create(church).Error
}

func (r *repository) GetChurchById(db *gorm.DB, id uint) (*Church, error) {
	var church Church 
	if err := db.First(&church, id).Error; err != nil {
		return nil, err
	}
	return &church, nil
}