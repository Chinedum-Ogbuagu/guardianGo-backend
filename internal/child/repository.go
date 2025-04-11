package child

import "gorm.io/gorm"


type Repository interface {
	Create(db *gorm.DB, child *Child) error 
	GetByID(db *gorm.DB, id uint) (*Child, error)
} 

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, child *Child) error {
	return db.Create(child).Error
}

func (r *repository) GetByID(db *gorm.DB, id uint) (*Child, error) {
	var c Child
	if err := db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}