package dropoff

import (
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
	GetDropSessionByCode(db *gorm.DB, code string) (*DropSession, error)
	GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error)
	GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error)
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
		guardianEntity = &guardian.Guardian{
			Name: req.Guardian.Name,
			Phone: req.Guardian.Phone,
			CreatedAt: time.Now(),
		}
		if err := s.guardianRepo.Create(tx, guardianEntity); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	dropSession := DropSession{
		UniqueCode: generateCode(),
		GuardianID: guardianEntity.ID,
		ChurchID: req.ChurchID,
		Note: req.Note,
		CreatedAt: time.Now(),
	}

	if err := s.dropRepo.CreateDropSession(tx, &dropSession); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, childInput := range req.Children {
	childEntity, err := s.childRepo.FindOrCreateChild(tx, childInput.Name, guardianEntity.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	dropOff := DropOff{
		DropSessionID: dropSession.ID,
		ChildID:       childEntity.ID,
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

func (s *service) GetDropSessionByCode(db *gorm.DB, code string) (*DropSession, error) {
	return s.dropRepo.GetDropSessionByCode(db, code)
}

func (s *service) GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error) {
	return s.dropRepo.GetDropOffsBySessionID(db, sessionID)
}

func (s *service) GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error) {
	return s.dropRepo.GetDropOffByID(db, id)
}