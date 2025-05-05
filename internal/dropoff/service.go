package dropoff

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/child"
	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/guardian"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateCode() string {
	return strconv.Itoa(100000 + rng.Intn(899999))
}

type GuardianInput struct {
	Name string `json:"name"`
	Phone string `json:"phone"`
}

type ChildInput struct {
	Name string `json:"name"`
	Class string `json:"class"`
	Bag bool `json:"bag"`
	Note string `json:"note"`
}


type Service interface {
	CreateDropSession(db *gorm.DB, req CreateDropSessionRequest) (*DropSession, error)
	GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error)
	GetDropSessionByCode(db *gorm.DB, code string) ([]*DropSession, error)
	GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error)
	GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error)
	GetDropSessionsByDate(db *gorm.DB, date time.Time, pagination Pagination) ([]DropSession, int64, error)
	MarkDropSessionPickedUp(db *gorm.DB, sessionID uint) error
}

type service struct {
	dropRepo Repository
	guardianRepo guardian.Repository
	childRepo child.Repository
}

func NewService(dr Repository, gr guardian.Repository, cr child.Repository) Service {
	return &service{
		dropRepo: dr,
		guardianRepo: gr,
		childRepo: cr,
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
		return nil, errors.New("this guardian has already dropped off children today")
	}


	dropSession := DropSession{
		UniqueCode: generateCode(),
		GuardianPhone: guardianEntity.Phone,
		GuardianName: guardianEntity.Name,
		GuardianID: guardianEntity.ID,
		ChurchID: req.ChurchID,
		Note: req.Note,
		CreatedAt: time.Now().In(loc),
	}

	if err := s.dropRepo.CreateDropSession(tx, &dropSession); err != nil {
		tx.Rollback()
		return nil, err
	}

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
}



	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.dropRepo.GetDropSessionByID(db, dropSession.ID)
}

func (s *service) GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error) {
	return s.dropRepo.GetDropSessionByID(db, id)
}

func (s *service) GetDropSessionByCode(db *gorm.DB, code string) ([]*DropSession, error) {
	return s.dropRepo.GetDropSessionByCode(db, code)
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