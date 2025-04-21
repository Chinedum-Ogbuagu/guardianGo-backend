package child

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindOrCreateChild(db *gorm.DB, name string, class string, guardianID uint) (*Child, error)
	GetChildrenByGuardian(db *gorm.DB, guardianID uint) ([]Child, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindOrCreateChild(db *gorm.DB, name string, class string, guardianID uint) (*Child, error) {
	var child Child
	err := db.Where("name = ? AND guardian_id = ?", name, guardianID).First(&child).Error
	if err == nil {
		return &child, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	child = Child{
		Name:       name,
		GuardianID: guardianID,
		Class:  class,
		CreatedAt:  time.Now(),
	}
	if err := db.Create(&child).Error; err != nil {
		return nil, err
	}

	return &child, nil
}

func (r *repository) GetChildrenByGuardian(db *gorm.DB, guardianID uint) ([]Child, error) {
	var children []Child
	if err := db.Where("guardian_id = ?", guardianID).Find(&children).Error; err != nil {
		return nil, err
	}
	return children, nil
}
