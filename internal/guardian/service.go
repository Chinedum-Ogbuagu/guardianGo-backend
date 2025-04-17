package guardian

import (
	"errors"

	"gorm.io/gorm"
)

type Service interface {
	FindOrCreateGuardian(db *gorm.DB, name, phone string) (*Guardian, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) FindOrCreateGuardian(db *gorm.DB, name, phone string) (*Guardian, error) {
	guardian, err := s.repo.FindByPhone(db, phone)
	if err == nil {
		return guardian, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newGuardian := &Guardian{
			Name:        name,
			Phone: phone,
		}
		if err := s.repo.Create(db, newGuardian); err != nil {
			return nil, err
		}
		return newGuardian, nil
	}
	return nil, err
}
