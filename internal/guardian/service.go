package guardian

import "gorm.io/gorm"

type Service interface {
	Create(db *gorm.DB, parent *Guardian) error
	GetByPhoneNumber(db *gorm.DB, phone  string) (*Guardian , error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(db *gorm.DB, guardian *Guardian) error {
	return s.repo.Create(db,guardian)

}

func (s *service) GetByPhoneNumber (db *gorm.DB, phone string ) (*Guardian, error) {
	return s.repo.GetByPhoneNumber(db, phone)
}