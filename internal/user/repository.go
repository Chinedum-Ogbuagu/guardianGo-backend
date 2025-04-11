package user

import "gorm.io/gorm"

type Repository interface {
	Create(db *gorm.DB, u *User) error
	GetByPhone(db *gorm.DB, phone string) (*User, error)
	ListByChurch(db *gorm.DB, churchID uint) ([]User, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(db *gorm.DB, u *User) error {
	return db.Create(u).Error
}

func (r *repository) GetByPhone(db *gorm.DB, phone string) (*User, error) {
	var u User
	if err := db.Where("phone_number = ?", phone).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) ListByChurch(db *gorm.DB, churchID uint) ([]User, error) {
	var users []User
	if err := db.Where("church_id = ?", churchID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
