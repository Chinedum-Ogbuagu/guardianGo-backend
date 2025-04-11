package parent

import "gorm.io/gorm"

type Service interface {
	Create(db *gorm.DB, parent *Parent) error
	GetByPhoneNumber(db *gorm.DB, phone  string) (*Parent , error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(db *gorm.DB, parent *Parent) error {
	return s.repo.Create(db,parent)

}

func (s *service) GetByPhoneNumber (db *gorm.DB, phone string ) (*Parent, error) {
	return s.repo.GetByPhoneNumber(db, phone)
}