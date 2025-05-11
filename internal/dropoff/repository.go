package dropoff

import (
	"time"

	"gorm.io/gorm"
)


type Repository interface {
	// Drop Session operations
	CreateDropSession(db *gorm.DB, ds *DropSession) error
	GetDropSessionByID(db *gorm.DB, id uint) (*DropSession, error)
	GetDropSessionByCode(db *gorm.DB, date time.Time, code string) ([]*DropSession, error)
	GetDropSessionsByDate(db *gorm.DB, date time.Time, pagination Pagination) ([]DropSession, int64, error)
	CheckGuardianDropSessionExistsForDate(db *gorm.DB, guardianID uint, date time.Time) (bool, error)
	UpdatePickupStatus(db *gorm.DB, sessionID uint, status string) error
	UpdateDropSessionImageURL(db *gorm.DB, sessionID string, imageURL string) error



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

func (r *repository) GetDropSessionByCode(db *gorm.DB, date time.Time, code string) ([]*DropSession, error) {
	var dropSessions []*DropSession

	loc, _ := time.LoadLocation("Africa/Lagos")
	date = date.In(loc)
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)

	if err := db.Preload("DropOffs").
		Where("unique_code LIKE ? AND created_at >= ? AND created_at < ?", "%"+code+"%", start, end).
		Find(&dropSessions).Error; err != nil {
		return nil, err
	}
	return dropSessions, nil
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
func (r *repository) UpdateDropSessionImageURL(db *gorm.DB, sessionID string, photoURL string) error {
	return db.Model(&DropSession{}).
		Where("unique_code = ?", sessionID).
		Update("photo_url", photoURL).Error
}


func (r *repository) GetDropSessionsByDate(db *gorm.DB, date time.Time, pagination Pagination) ([]DropSession, int64, error) {
	var sessions []DropSession
	var totalCount int64
	loc, _ := time.LoadLocation("Africa/Lagos")
	date = date.In(loc)
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)

	query := db.Preload("DropOffs").
		Where("created_at >= ? AND created_at < ?", start, end)

	// Conditionally apply pagination if parameters are provided
	if pagination.Page > 0 && pagination.PageSize > 0 {
		if err := db.Model(&DropSession{}).
			Where("created_at >= ? AND created_at < ?", start, end).
			Count(&totalCount).Error; err != nil {
			return nil, 0, err
		}
		offset := (pagination.Page - 1) * pagination.PageSize
		query = query.Offset(offset).Limit(pagination.PageSize)
	} else {
		// If no pagination parameters, get all records
		if err := db.Model(&DropSession{}).
			Where("created_at >= ? AND created_at < ?", start, end).
			Count(&totalCount).Error; err != nil {
			return nil, 0, err
		}
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, totalCount, nil
}


func (r *repository) GetDropOffsBySessionID(db *gorm.DB, sessionID uint) ([]DropOff, error) {
	var dropOffs []DropOff
	if err := db.Where("drop_session_id = ?", sessionID).Find(&dropOffs).Error; err != nil {
		return nil, err
	}
	return dropOffs, nil
}


func (r *repository) CheckGuardianDropSessionExistsForDate(db *gorm.DB, guardianID uint, date time.Time) (bool, error) {
	var count int64
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	if err := db.Model(&DropSession{}).
		Where("guardian_id = ? AND created_at >= ? AND created_at < ?", guardianID, start, end).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) UpdatePickupStatus(db *gorm.DB, sessionID uint, status string) error {
	return db.Model(&DropSession{}).Where("id = ?", sessionID).Update("pickup_status", status).Error
}
