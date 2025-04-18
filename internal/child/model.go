package child

import "time"

type Child struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	Class      string    `json:"class"`           // NEW: child class (e.g. Nursery, Toddlers)
	Note       string    `json:"note"`            // NEW: optional medical or personal notes
	BagStatus  bool      `json:"bag_status"`
	Age        int       `json:"age"`
	GuardianID uint      `json:"guardian_id"`
	CreatedAt  time.Time `json:"created_at"`
}