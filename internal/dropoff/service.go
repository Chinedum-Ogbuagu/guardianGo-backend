package dropoff

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/google/uuid"

	"time"

	"gorm.io/gorm"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/child"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/guardian"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/messaging"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/otp"
)

func generatePublicCode() string {
	return uuid.New().String()[:6] // e.g., "a1b2c3"
}

var letters = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ23456789") // excludes confusing chars like I, O, 1, 0

func generatePickupSecret() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

type GuardianInput struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type ChildInput struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Bag   bool   `json:"bag"`
	Note  string `json:"note"`
}

type Service interface {
	CreateDropSession(db *gorm.DB, req CreateDropSessionRequest) (*DropSession, error)
	GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error)
	GetDropSessionByCode(db *gorm.DB, date time.Time, code string) ([]*DropSession, error)
	GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error)
	GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error)
	VerifyPickupCode(db *gorm.DB, code string, secret string) (*DropSession, error)

	GetDropSessionsByDate(db *gorm.DB, date time.Time, pagination Pagination) ([]DropSession, int64, error)
	MarkDropSessionPickedUp(db *gorm.DB, sessionID uint) error
	UpdateDropSessionImageURL(db *gorm.DB, sessionID string, photoURL string) error
}

type service struct {
	dropRepo         Repository
	guardianRepo     guardian.Repository
	childRepo        child.Repository
	messagingService messaging.Service
	otpService       otp.Service
}

func NewService(dr Repository, gr guardian.Repository, cr child.Repository, ms messaging.Service, os otp.Service) Service {
	return &service{
		dropRepo:         dr,
		guardianRepo:     gr,
		childRepo:        cr,
		messagingService: ms,
		otpService:       os,
	}
}

func (s *service) CreateDropSession(db *gorm.DB, req CreateDropSessionRequest) (*DropSession, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	guardianEntity, err := s.guardianRepo.FindByPhone(tx, req.Guardian.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			guardianEntity = &guardian.Guardian{
				Name:      req.Guardian.Name,
				Phone:     req.Guardian.Phone,
				Email:     req.Guardian.Email,
				CreatedAt: time.Now(),
			}
			if err := s.guardianRepo.Create(tx, guardianEntity); err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			tx.Rollback()
			return nil, err
		}
	}
	loc, _ := time.LoadLocation("Africa/Lagos")
	alreadyDroppedToday, err := s.dropRepo.CheckGuardianDropSessionExistsForDate(tx, guardianEntity.ID, time.Now().In(loc))
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if alreadyDroppedToday {
		tx.Rollback()
		return nil, errors.New("this guardian has already dropped off children today, You can Still Add Children to this session")
	}
	secret := generatePickupSecret()

	dropSession := DropSession{
		UniqueCode:    generatePublicCode(),
		PickupSecret:  secret,
		GuardianPhone: guardianEntity.Phone,
		GuardianName:  guardianEntity.Name,
		GuardianID:    guardianEntity.ID,
		ChurchID:      req.ChurchID,
		Note:          req.Note,
		CreatedAt:     time.Now().In(loc),
	}

	if err := s.dropRepo.CreateDropSession(tx, &dropSession); err != nil {
		tx.Rollback()
		return nil, err
	}

	var childNames []string
	for _, childInput := range req.Children {
		childEntity, err := s.childRepo.FindOrCreateChild(tx, childInput.Name, childInput.Class, guardianEntity.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		dropOff := DropOff{
			DropSessionID: dropSession.ID,
			ChildID:       childEntity.ID,
			ChildName:     childEntity.Name,
			Class:         childInput.Class,
			BagStatus:     childInput.Bag,
			Note:          childInput.Note,
			DropOffTime:   time.Now(),
			CreatedAt:     time.Now(),
		}

		if err := s.dropRepo.CreateDropOff(tx, &dropOff); err != nil {
			tx.Rollback()
			return nil, err
		}
		childNames = append(childNames, childEntity.Name)

	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	// Send the pickup secret via email
	_, err = s.otpService.SendEmailOTP(guardianEntity.Email, secret)
	fmt.Printf("Error sending pickup secret email: %s", guardianEntity.Email)
	if err != nil {
		fmt.Printf("Error sending pickup secret email: %v\n", err)
		// Consider if this error should be fatal or just logged.
		// For now, we'll log it and continue with the SMS.
	}

	whatsappData := map[string]string{
		"pin": secret,
	}
	whatsappResponse, err := s.otpService.SendWhatsAppOTP(guardianEntity.Phone, whatsappData)
	fmt.Printf("Attempting to send WhatsApp OTP to: %s\n", guardianEntity.Phone)
	if err != nil {
		fmt.Printf("Error sending WhatsApp OTP: %v\n", err)
		// Handle WhatsApp sending errors (log, maybe retry, etc.)
	} else {
		fmt.Printf("WhatsApp OTP sent successfully: %+v\n", whatsappResponse)
	}
	go s.messagingService.SendDropSessionEmail(messaging.DropSessionEmailPayload{
		GuardianName:  guardianEntity.Name,
		GuardianEmail: guardianEntity.Email,
		Children:      childNames,
		Secret:        secret,
		Date:          time.Now().In(loc).Format("Monday, Jan 2 2006"),
		ChurchName:    "Your Church Name", // optional: pull from church repo
	})
	return s.dropRepo.GetDropSessionByID(db, dropSession.ID)
}
func (s *service) VerifyPickupCode(db *gorm.DB, code string, secret string) (*DropSession, error) {
	return s.dropRepo.VerifyPickupSecret(db, code, secret)
}

func (s *service) GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error) {
	return s.dropRepo.GetDropSessionByID(db, id)
}

func (s *service) GetDropSessionByCode(db *gorm.DB, date time.Time, code string) ([]*DropSession, error) {
	return s.dropRepo.GetDropSessionByCode(db, date, code)
}

func (s *service) GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error) {
	return s.dropRepo.GetDropOffsBySessionID(db, sessionID)
}

func (s *service) GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error) {
	return s.dropRepo.GetDropOffByID(db, id)
}
func (s *service) GetDropSessionsByDate(db *gorm.DB, date time.Time, pagination Pagination) ([]DropSession, int64, error) {
	return s.dropRepo.GetDropSessionsByDate(db, date, pagination)
}

func (s *service) MarkDropSessionPickedUp(db *gorm.DB, sessionID uint) error {
	return db.Model(&DropSession{}).
		Where("id = ?", sessionID).
		Update("pickup_status", "completed").Error
}
func (s *service) UpdateDropSessionImageURL(db *gorm.DB, sessionID string, photoURL string) error {
	println("Updating photo URL for session ID:", sessionID, "to", photoURL)
	return s.dropRepo.UpdateDropSessionImageURL(db, sessionID, photoURL)
}
