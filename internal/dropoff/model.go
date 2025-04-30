package dropoff

import (
	"time"

	"github.com/gofrs/uuid"
)

type DropSession struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UniqueCode string    `json:"unique_code"`
	GuardianID uint      `json:"guardian_id"`
	GuardianPhone string  `json:"guardian_phone"`
	GuardianName string  `json:"guardian_name"`
	ChurchID  *uuid.UUID `json:"church_id"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	PickupStatus     string    `gorm:"default:'awaiting'" json:"pickup_status"` 
	DropOffs   []DropOff `json:"drop_offs,omitempty" gorm:"foreignKey:DropSessionID"`
}


type DropOff struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DropSessionID uint      `json:"drop_session_id"`
	ChildID       uint      `json:"child_id"`
	Class         string    `json:"class"`
	ChildName	 string    `json:"child_name"`
	BagStatus     bool      `json:"bag_status"`
	Note          string    `json:"note"`
	DropOffTime   time.Time  `json:"drop_off_time"`
	CreatedAt     time.Time `json:"created_at"`
}