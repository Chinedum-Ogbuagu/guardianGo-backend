package user

import "gorm.io/gorm"

type Repository interface {
	Create(db *gorm.DB, user *User) error
	FindByPhone(db *gorm.DB, phone string) (*User, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func (r *repository) FindByPhone(db *gorm.DB, phone string) (*User, error) {
	var user User
	if err := db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
