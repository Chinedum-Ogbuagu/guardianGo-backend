package pickup

import "gorm.io/gorm"

type Repository interface {
	// PickupSession operations
	CreatePickupSession(db *gorm.DB, ps *PickupSession) error
	GetPickupSessionByID(db *gorm.DB, id uint) (*PickupSession, error)
	GetPickupSessionByDropSessionID(db *gorm.DB, dropSessionID uint) (*PickupSession, error)
	GetPickupSessionByCode(db *gorm.DB, code string) (*PickupSession, error)
	
	// Pickup operations
	CreatePickup(db *gorm.DB, p *Pickup) error
	GetPickupByID(db *gorm.DB, id uint) (*Pickup, error)
	GetPickupsBySessionID(db *gorm.DB, pickupSessionID uint) ([]Pickup, error)
	GetPickupByChildAndDropSessionID(db *gorm.DB, childID, dropSessionID uint) (*Pickup, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreatePickupSession(db *gorm.DB, ps *PickupSession) error {
	return db.Create(ps).Error
}

func (r *repository) GetPickupSessionByID(db *gorm.DB, id uint) (*PickupSession, error) {
	var ps PickupSession
	if err := db.Preload("Pickups").First(&ps, id).Error; err != nil {
		return nil, err
	}
	return &ps, nil
}

func (r *repository) GetPickupSessionByDropSessionID(db *gorm.DB, dropSessionID uint) (*PickupSession, error) {
	var ps PickupSession
	if err := db.Preload("Pickups").Where("drop_session_id = ?", dropSessionID).First(&ps).Error; err != nil {
		return nil, err
	}
	return &ps, nil
}

func (r *repository) GetPickupSessionByCode(db *gorm.DB, code string) (*PickupSession, error) {
	var ps PickupSession
	if err := db.Preload("Pickups").Where("unique_code = ?", code).First(&ps).Error; err != nil {
		return nil, err
	}
	return &ps, nil
}

func (r *repository) CreatePickup(db *gorm.DB, p *Pickup) error {
	return db.Create(p).Error
}

func (r *repository) GetPickupByID(db *gorm.DB, id uint) (*Pickup, error) {
	var p Pickup
	if err := db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) GetPickupsBySessionID(db *gorm.DB, pickupSessionID uint) ([]Pickup, error) {
	var pickups []Pickup
	if err := db.Where("pickup_session_id = ?", pickupSessionID).Find(&pickups).Error; err != nil {
		return nil, err
	}
	return pickups, nil
}

func (r *repository) GetPickupByChildAndDropSessionID(db *gorm.DB, childID, dropSessionID uint) (*Pickup, error) {
	var p Pickup
	if err := db.Joins("JOIN pickup_sessions ON pickups.pickup_session_id = pickup_sessions.id").
		Where("pickups.child_id = ? AND pickup_sessions.drop_session_id = ?", childID, dropSessionID).
		First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
