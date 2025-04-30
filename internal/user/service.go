package user

import (
	"errors"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Service interface {
	RegisterUser(db *gorm.DB, name, phone string, role Role, churchID uuid.UUID) (*User, error)
	GetUserByPhone(db *gorm.DB, phone string) (*User, error)
	FindOrCreateUserByPhone(db *gorm.DB, phone string, name string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterUser(db *gorm.DB, name, phone string, role Role, churchID uuid.UUID) (*User, error) {
	existing, err := s.repo.FindByPhone(db, phone)
	if err == nil {
		return existing, nil // Return existing user
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &User{
		Name:     name,
		Phone:    phone,
		Role:     role,
		ChurchID: &churchID,
	}
	if err := s.repo.Create(db, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetUserByPhone(db *gorm.DB, phone string) (*User, error) {
	return s.repo.FindByPhone(db, phone)
}

func (s *service) FindOrCreateUserByPhone(db *gorm.DB, phone, name string ) (*User, error) {
	user, err := s.repo.FindByPhone(db, phone)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Default to RoleAttendant and ChurchID 0 for now (can adjust later)
		newUser := &User{
			Name:     name,
			Phone:    phone,
			Role:     RoleAttendant,
			ChurchID: nil, // Set to nil if no ChurchID is provided
		}
		if err := s.repo.Create(db, newUser); err != nil {
			return nil, err
		}
		return newUser, nil
	}
	return nil, err
}
