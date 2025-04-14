package pickup

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/dropoff"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateCode() string {
	return strconv.Itoa(100000 + rng.Intn(899999))
}

type Service interface {
	// PickupSession operations
	CreatePickupSession(db *gorm.DB, dropSessionID, guardianID, verifiedByID uint, notes string) (*PickupSession, error)
	GetPickupSessionByID(db *gorm.DB, id uint) (*PickupSession, error)
	GetPickupSessionByDropSessionID(db *gorm.DB, dropSessionID uint) (*PickupSession, error)
	ValidatePickupCode(db *gorm.DB, dropSessionID uint, code string) (bool, error)
	
	// Pickup operations
	GetPickupByChildAndDropSessionID(db *gorm.DB, childID, dropSessionID uint) (*Pickup, error)
}

type service struct {
	repo Repository
	dropoffService dropoff.Service // Dependency on dropoff service
}

// Include a reference to the dropoff service for cross-domain validation
func NewService(r Repository, ds dropoff.Service) Service {
	return &service{
		repo: r,
		dropoffService: ds,
	}
}

func (s *service) CreatePickupSession(db *gorm.DB, dropSessionID, guardianID, verifiedByID uint, notes string) (*PickupSession, error) {
	// Check if pickup session already exists for this drop session
	existingSession, err := s.repo.GetPickupSessionByDropSessionID(db, dropSessionID)
	if err == nil && existingSession != nil {
		return nil, errors.New("pickup session already exists for this drop session")
	}
	
	// Get drop session to validate it exists and get all children
	if _, err := s.dropoffService.GetDropSessionByID(db, dropSessionID); err != nil {
		return nil, errors.New("drop session not found")
	}
	
	// Get all drop-offs for this session
	dropOffs, err := s.dropoffService.GetDropOffsBySessionID(db, dropSessionID)
	if err != nil {
		return nil, errors.New("failed to retrieve drop-offs")
	}
	
	if len(dropOffs) == 0 {
		return nil, errors.New("no children to pick up in this drop session")
	}
	
	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// Create pickup session
	pickupSession := PickupSession{
		DropSessionID: dropSessionID,
		GuardianID:    guardianID,
		VerifiedByID:  verifiedByID,
		UniqueCode:    generateCode(),
		VerifiedAt:    time.Now(),
		Notes:         notes,
		CreatedAt:     time.Now(),
	}
	
	if err := s.repo.CreatePickupSession(tx, &pickupSession); err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// Create individual pickups for each child in the drop session
	
	for _, dropOff := range dropOffs {
		pickup := Pickup{
			PickupSessionID: pickupSession.ID,
			ChildID:        dropOff.ChildID,
			DropOffID:      dropOff.ID,
			PickupTime:     time.Now(),
			CreatedAt:      time.Now(),
		}
		
		if err := s.repo.CreatePickup(tx, &pickup); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	
	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	
	// Return the complete pickup session with all pickups
	return s.repo.GetPickupSessionByID(db, pickupSession.ID)
}

func (s *service) GetPickupSessionByID(db *gorm.DB, id uint) (*PickupSession, error) {
	return s.repo.GetPickupSessionByID(db, id)
}

func (s *service) GetPickupSessionByDropSessionID(db *gorm.DB, dropSessionID uint) (*PickupSession, error) {
	return s.repo.GetPickupSessionByDropSessionID(db, dropSessionID)
}

func (s *service) ValidatePickupCode(db *gorm.DB, dropSessionID uint, code string) (bool, error) {
	// First check if drop session exists
	_, err := s.dropoffService.GetDropSessionByCode(db, code)
	if err != nil {
		return false, errors.New("invalid drop session code")
	}
	
	return true, nil
}

func (s *service) GetPickupByChildAndDropSessionID(db *gorm.DB, childID, dropSessionID uint) (*Pickup, error) {
	return s.repo.GetPickupByChildAndDropSessionID(db, childID, dropSessionID)
}