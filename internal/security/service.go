package security

import (
	"time"

	"gorm.io/gorm"
)

type Service interface {
	FlagIssue(db *gorm.DB , flag *SecurityFlag) error
	GetAllFlags(db *gorm.DB) ([]SecurityFlag, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) FlagIssue(db *gorm.DB, flag *SecurityFlag) error {
	flag.CreatedAt = time.Now().Format(time.RFC3339)
	return s.repo.Create(db, flag)
}

func (s *service) GetAllFlags(db *gorm.DB) ([]SecurityFlag, error) {
	flags, err := s.repo.ListAll(db)
	if err != nil {
		return nil, err
	}
	return flags, nil
}