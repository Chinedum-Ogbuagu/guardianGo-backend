package auth

import (
	"errors"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/otp"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/user"
	"gorm.io/gorm"
)

type Service interface {
	RequestOTP(db *gorm.DB, phone string, name string, dropOffID uint) error
	VerifyOTPAndLogin(db *gorm.DB, phone, otpCode string, name string, purpose string) (*user.User, string, error)
}

type service struct {
	otpService  otp.Service
	userService user.Service
}

func NewService(otpService otp.Service, userService user.Service) Service {
	return &service{
		otpService:  otpService,
		userService: userService,
	}
}

func (s *service) RequestOTP(db *gorm.DB, phone string, name string, dropOffID uint) error {
	_, err := s.otpService.SendOTP(db, phone, "login", dropOffID)
	return err
}

func (s *service) VerifyOTPAndLogin(db *gorm.DB, phone, otpCode string, name string, purpose string) (*user.User, string, error) {
	isValid, err := s.otpService.VerifyOTP(db, phone, otpCode, purpose)
	if err != nil || !isValid {
		return nil, "", errors.New("invalid or expired OTP")
	}

	foundUser, err := s.userService.GetUserByPhone(db, phone)
	if err != nil {
		return nil, "", errors.New("user not found")
	}

	token, err := GenerateJWT(foundUser.ID, foundUser.Role)
	if err != nil {
		return nil, "", err
	}

	return foundUser, token, nil
}

