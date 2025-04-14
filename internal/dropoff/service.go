package dropoff

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateCode() string {
	return strconv.Itoa(100000 + rng.Intn(899999))
}

type Service interface {
	
	CreateDropSession(db *gorm.DB, ds *DropSession, childrenIDs []uint, classes []string, bagStatuses []bool) (*DropSession, error)
	GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error)
	GetDropSessionByCode(db *gorm.DB, code string) (*DropSession, error)
	
	
	AddChildToSession(db *gorm.DB, sessionID uint, childID uint, class string, bagStatus bool, note string) (*DropOff, error)
	GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error)
	GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateDropSession(db *gorm.DB, ds *DropSession, childrenIDs []uint, classes []string, bagStatuses []bool) (*DropSession, error) {
	
	if len(childrenIDs) == 0 {
		return nil, errors.New("at least one child must be provided")
	}
	
	if len(childrenIDs) != len(classes) || len(childrenIDs) != len(bagStatuses) {
		return nil, errors.New("mismatch between childrenIDs, classes, and bagStatuses arrays")
	}
	
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	
	ds.UniqueCode = generateCode()
	ds.CreatedAt = time.Now()
	

	if err := s.repo.CreateDropSession(tx, ds); err != nil {
		tx.Rollback()
		return nil, err
	}
	
	
	currentTime := time.Now()
	for i, childID := range childrenIDs {
		dropOff := DropOff{
			DropSessionID: ds.ID,
			ChildID:      childID,
			Class:        classes[i],
			BagStatus:    bagStatuses[i],
			DropOffTime:  currentTime,
			CreatedAt:    time.Now(),
		}
		
		if err := s.repo.CreateDropOff(tx, &dropOff); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	
	
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	
	
	return s.repo.GetDropSessionByID(db, ds.ID)
}

func (s *service) GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error) {
	return s.repo.GetDropSessionByID(db, id)
}

func (s *service) GetDropSessionByCode(db *gorm.DB, code string) (*DropSession, error) {
	return s.repo.GetDropSessionByCode(db, code)
}

func (s *service) AddChildToSession(db *gorm.DB, sessionID uint, childID uint, class string, bagStatus bool, note string) (*DropOff, error) {
	
	_, err := s.repo.GetDropSessionByID(db, sessionID)
	if err != nil {
		return nil, errors.New("drop session not found")
	}
	
	dropOff := DropOff{
		DropSessionID: sessionID,
		ChildID:      childID,
		Class:        class,
		BagStatus:    bagStatus,
		Note:         note,
		DropOffTime:  time.Now(),
		CreatedAt:    time.Now(),
	}
	
	if err := s.repo.CreateDropOff(db, &dropOff); err != nil {
		return nil, err
	}
	
	return &dropOff, nil
}

func (s *service) GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error) {
	return s.repo.GetDropOffsBySessionID(db, sessionID)
}

func (s *service) GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error) {
	return s.repo.GetDropOffByID(db, id)
}