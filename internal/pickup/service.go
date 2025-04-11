package pickup

import (
	"time"

	"gorm.io/gorm"
)

type Service interface {
	Create(db *gorm.DB, p *Pickup) error
	GetByID(db *gorm.DB, id uint) (*Pickup, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(db *gorm.DB, p *Pickup) error {
	p.PickupTime = time.Now().Format(time.RFC3339)
	return s.repo.Create(db, p)
}

func (s *service) GetByID(db *gorm.DB, id uint) (*Pickup, error) {
	return s.repo.GetByID(db, id)
}