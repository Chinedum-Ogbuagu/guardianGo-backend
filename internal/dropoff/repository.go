package dropoff

import (
	"time"

	"gorm.io/gorm"
)


type Repository interface {
	// Drop Session operations
	CreateDropSession(db *gorm.DB, ds *DropSession) error
	GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error)
	GetDropSessionByCode(db *gorm.DB, code string) (*DropSession, error)
	GetDropSessionsByDate(db *gorm.DB, date time.Time) ([]DropSession, error)

	// Drop Off operations
	CreateDropOff(db *gorm.DB, d *DropOff) error
	GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error)
	GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateDropSession(db *gorm.DB, ds *DropSession) error {
	return db.Create(ds).Error
}

func (r *repository) GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error) {
	var ds DropSession
	if err := db.Preload("DropOffs").First(&ds, id).Error; err != nil {
		return nil, err
	}
	return &ds, nil
}

func (r *repository) GetDropSessionByCode(db *gorm.DB, code string) (*DropSession, error) {
	var ds DropSession
	if err := db.Preload("DropOffs").Where("unique_code = ?", code).First(&ds).Error; err != nil {
		return nil, err
	}
	return &ds, nil
}

func (r *repository) CreateDropOff(db *gorm.DB, d *DropOff) error {
	return db.Create(d).Error
}

func (r *repository) GetDropOffByID(db *gorm.DB, id uint) (*DropOff, error) {
	var d DropOff
	if err := db.First(&d, id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repository) GetDropSessionsByDate(db *gorm.DB, date time.Time) ([]DropSession, error) {
	var sessions []DropSession

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	if err := db.Preload("DropOffs").
		Where("created_at >= ? AND created_at < ?", start, end).
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}


func (r *repository) GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error) {
	var dropOffs []DropOff
	if err := db.Where("drop_session_id = ?", sessionID).Find(&dropOffs).Error; err != nil {
		return nil, err
	}
	return dropOffs, nil
}
