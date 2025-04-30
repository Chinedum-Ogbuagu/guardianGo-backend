package church

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(db *gorm.DB, church *Church) error
	GetByID(db *gorm.DB, id uuid.UUID) (*Church, error)
	GetAll(db *gorm.DB) ([]*Church, error)
	Update(db *gorm.DB, church *Church) error
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

func (s *service ) GetByID(db *gorm.DB, id uuid.UUID) (*Church, error) {
	return s.repo.GetChurchById(db, id)
}
func (s *service) GetAll(db *gorm.DB) ([]*Church, error) {
	return s.repo.GetAllChurches(db)
}

func (s *service) Update(db *gorm.DB, church *Church) error {
	return s.repo.UpdateChurch(db, church)
}