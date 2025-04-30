package church

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateChurch(db *gorm.DB, church *Church) error
	GetChurchById(db *gorm.DB, id uuid.UUID) (*Church, error)
	GetAllChurches(db *gorm.DB) ([]*Church, error)
	UpdateChurch(db *gorm.DB, church *Church) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateChurch(db *gorm.DB, church *Church) error {
	return db.Create(church).Error
}

func (r *repository) GetChurchById(db *gorm.DB, id uuid.UUID) (*Church, error) {
	var church Church 
	if err := db.First(&church, id).Error; err != nil {
		return nil, err
	}
	return &church, nil
}
func (r *repository) GetAllChurches(db *gorm.DB) ([]*Church, error) {
	var churches []*Church
	if err := db.Find(&churches).Error; err != nil {
		return nil, err
	}
	return churches, nil
}
func (r *repository) UpdateChurch(db *gorm.DB, church *Church) error {
	return db.Save(church).Error
}