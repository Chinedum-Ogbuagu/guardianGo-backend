package pickup

import "gorm.io/gorm"

type Repository interface {
	CreatePickupSession(db *gorm.DB, session *PickupSession) error
	CreatePickup(db *gorm.DB, pickup *Pickup) error
	GetPickupSessionByDropSessionID(db *gorm.DB, dropSessionID uint) (*PickupSession, error)
	GetPickupsByPickupSessionID(db *gorm.DB, pickupSessionID uint) ([]Pickup, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreatePickupSession(db *gorm.DB, ps *PickupSession) error {
	return db.Create(ps).Error
}

func (r *repository) CreatePickup(db *gorm.DB, p *Pickup) error {
	return db.Create(p).Error
}

func (r *repository) GetPickupSessionByDropSessionID(db *gorm.DB, dropSessionID uint) (*PickupSession, error) {
	var ps PickupSession
	if err := db.Preload("Pickups").Where("drop_session_id = ?", dropSessionID).First(&ps).Error; err != nil {
		return nil, err
	}
	return &ps, nil
}
func (r *repository) GetPickupsByPickupSessionID(db *gorm.DB, pickupSessionID uint) ([]Pickup, error) {
	var pickups []Pickup
	err := db.Where("pickup_session_id = ?", pickupSessionID).Find(&pickups).Error
	return pickups, err
}