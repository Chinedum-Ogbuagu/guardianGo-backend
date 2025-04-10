package church

import "gorm.io/gorm"

type Service interface {
	Create(db *gorm.DB, church *Church) error
	GetByID(db *gorm.DB, id uint) (*Church, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}

}

func (s *service) Create(db *gorm.DB, church *Church) error {
	return s.repo.CreateChurch(db, church)
} 

func (s *service ) GetByID(db *gorm.DB, id uint) (*Church, error) {
	return s.repo.GetChurchById(db, id)
}