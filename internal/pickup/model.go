package pickup

import (
	"time"
)


type PickupSession struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DropSessionID uint      `gorm:"uniqueIndex" json:"drop_session_id"`
	GuardianID    uint      `json:"guardian_id"`
	VerifiedByID  uint      `json:"verified_by_id"`
	UniqueCode    string    `json:"unique_code"`
	VerifiedAt    time.Time `json:"verified_at"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	Pickups       []Pickup  `json:"pickups,omitempty" gorm:"foreignKey:PickupSessionID"`
}


type Pickup struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PickupSessionID uint      `json:"pickup_session_id"`
	ChildID        uint      `json:"child_id"`
	DropOffID      uint      `json:"drop_off_id"`
	PickupTime     time.Time    `json:"pickup_time"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
}