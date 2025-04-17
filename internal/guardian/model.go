package guardian

import "time"

type Guardian struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	Name      string   `json:"name"`
	Phone     string   `gorm:"uniqueIndex" json:"phone_number"`
	CreatedAt time.Time `json:"created_at"`
}
