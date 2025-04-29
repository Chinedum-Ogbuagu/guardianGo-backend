package auth

import "gorm.io/gorm"

type Repository interface {
	CreateAuthSession(db *gorm.DB, session *AuthSession) error
	GetAuthSessionByToken(db *gorm.DB, token string) (*AuthSession, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateAuthSession(db *gorm.DB, session *AuthSession) error {
	return db.Create(session).Error
}

func (r *repository) GetAuthSessionByToken(db *gorm.DB, token string) (*AuthSession, error) {
	var session AuthSession
	if err := db.Where("token = ?", token).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
