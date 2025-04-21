package pickup

import (
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
	ConfirmPickup(db *gorm.DB, dropSession dropoff.DropSession, verifiedBy uint, notes string) (*PickupSession, error)
	GetPickupSessionByDropSessionID(db *gorm.DB, dropSessionID uint) (*PickupSession, error)
		ConfirmPickupSession(db *gorm.DB, dropSessionID, guardianID, verifiedByID uint, notes string) (*PickupSession, error)

}

type service struct {
	repo      Repository
	dropRepo  dropoff.Repository
}

func NewService(r Repository, dropRepo dropoff.Repository) Service {
	return &service{repo: r, dropRepo: dropRepo}
}

func (s *service) ConfirmPickup(db *gorm.DB, session dropoff.DropSession, verifiedBy uint, notes string) (*PickupSession, error) {
	tx := db.Begin()

	pickupSession := &PickupSession{
		DropSessionID: session.ID,
		GuardianID:    session.GuardianID,
		VerifiedByID:  verifiedBy,
		VerifiedAt:    time.Now(),
		CreatedAt:     time.Now(),
		UniqueCode:    generateCode(),
		Notes:         notes,
	}

	if err := s.repo.CreatePickupSession(tx, pickupSession); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, d := range session.DropOffs {
		p := &Pickup{
			PickupSessionID: pickupSession.ID,
			ChildID:         d.ChildID,
			DropOffID:       d.ID,
			PickupTime:      time.Now(),
			CreatedAt:       time.Now(),
		}
		if err := s.repo.CreatePickup(tx, p); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Update drop session status
	if err := tx.Model(&dropoff.DropSession{}).Where("id = ?", session.ID).Update("pickup_session_completed", true).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return pickupSession, nil
}

func (s *service) GetPickupSessionByDropSessionID(db *gorm.DB, dropSessionID uint) (*PickupSession, error) {
	return s.repo.GetPickupSessionByDropSessionID(db, dropSessionID)
}
func (s *service) ConfirmPickupSession(db *gorm.DB, dropSessionID, guardianID, verifiedByID uint, notes string) (*PickupSession, error) {
	tx := db.Begin()

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

	// Fetch all DropOffs tied to the DropSession
	var dropOffs []dropoff.DropOff
	if err := tx.Where("drop_session_id = ?", dropSessionID).Find(&dropOffs).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, do := range dropOffs {
		p := Pickup{
			PickupSessionID: pickupSession.ID,
			ChildID:         do.ChildID,
			DropOffID:       do.ID,
			PickupTime:      time.Now(),
			CreatedAt:       time.Now(),
		}
		if err := s.repo.CreatePickup(tx, &p); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Update DropSession status to completed
	if err := tx.Model(&dropoff.DropSession{}).Where("id = ?", dropSessionID).Update("pickup_session_completed", true).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &pickupSession, nil
}