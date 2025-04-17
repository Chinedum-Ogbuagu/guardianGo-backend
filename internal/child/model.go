package child

import "time"

type Child struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	GuardianID uint      `json:"guardian_id"`
	CreatedAt  time.Time `json:"created_at"`
	BagStatus  bool   `json:"bag_status"`
	Age        int    `json:"age"`
}
