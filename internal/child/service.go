package child

import "gorm.io/gorm"

type Service interface {
	FindOrCreateChild(db *gorm.DB, name string, guardianID uint) (*Child, error)
	GetChildrenByGuardian(db *gorm.DB, guardianID uint) ([]Child, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) FindOrCreateChild(db *gorm.DB, name string, guardianID uint) (*Child, error) {
	return s.repo.FindOrCreateChild(db, name, guardianID)
}

func (s *service) GetChildrenByGuardian(db *gorm.DB, guardianID uint) ([]Child, error) {
	return s.repo.GetChildrenByGuardian(db, guardianID)
}
