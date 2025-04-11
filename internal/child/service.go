package child

import "gorm.io/gorm"

type Service interface {
	Create(db *gorm.DB, child *Child) error
	GetByID(db *gorm.DB, id uint) (*Child, error)
}

type service struct {
	repo Repository

}


func NewService (r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(db *gorm.DB, child *Child) error {
	return s.repo.Create(db, child)
}

func (s *service) GetByID(db *gorm.DB, id uint) (*Child, error) {
	return s.repo.GetByID(db, id)
}
