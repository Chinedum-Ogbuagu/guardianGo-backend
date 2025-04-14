package dropoff

import (
	"time"
)

type DropSession struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UniqueCode string    `json:"unique_code"`
	GuardianID uint      `json:"guardian_id"`
	ChurchID   uint      `json:"church_id"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	DropOffs   []DropOff `json:"drop_offs,omitempty" gorm:"foreignKey:DropSessionID"`
}


type DropOff struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DropSessionID uint      `json:"drop_session_id"`
	ChildID       uint      `json:"child_id"`
	Class         string    `json:"class"`
	BagStatus     bool      `json:"bag_status"`
	Note          string    `json:"note"`
	DropOffTime   time.Time  `json:"drop_off_time"`
	CreatedAt     time.Time `json:"created_at"`
}