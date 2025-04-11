package user

import "gorm.io/gorm"

type Service interface {
	Create(db *gorm.DB, u *User) error
	GetByPhone(db *gorm.DB, phone string) (*User, error)
	ListByChurch(db *gorm.DB, churchID uint) ([]User, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(db *gorm.DB, u *User) error {
	return s.repo.Create(db, u)
}

func (s *service) GetByPhone(db *gorm.DB, phone string) (*User, error) {
	return s.repo.GetByPhone(db, phone)
}

func (s *service) ListByChurch(db *gorm.DB, churchID uint) ([]User, error) {
	return s.repo.ListByChurch(db, churchID)
}