package guardian

import (
	"errors"

	"gorm.io/gorm"
)

type Service interface {
	FindOrCreateGuardian(db *gorm.DB, name, phone, email string) (*Guardian, error)
	FindGuardianByPhone(db *gorm.DB, phone string) (*Guardian, error)
	GetChildrenByGuardianPhone(db *gorm.DB, phone string) ([]ChildInfo, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) FindOrCreateGuardian(db *gorm.DB, name, phone, email string) (*Guardian, error) {
	guardian, err := s.repo.FindByPhone(db, phone)
	if err == nil {
		return guardian, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newGuardian := &Guardian{
			Name:  name,
			Phone: phone,
			Email: email,
		}
		if err := s.repo.Create(db, newGuardian); err != nil {
			return nil, err
		}
		return newGuardian, nil
	}
	return nil, err
}
func (s *service) FindGuardianByPhone(db *gorm.DB, phone string) (*Guardian, error) {
	return s.repo.FindByPhone(db, phone)
}
func (s *service) GetChildrenByGuardianPhone(db *gorm.DB, phone string) ([]ChildInfo, error) {
	return s.repo.GetChildrenByGuardianPhone(db, phone)
}
