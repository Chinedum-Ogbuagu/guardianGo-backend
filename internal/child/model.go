package child

import "time"

type Child struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	Class      string    `json:"class"`     
	Age        int       `json:"age"`
	GuardianID uint      `json:"guardian_id"`
	CreatedAt  time.Time `json:"created_at"`
}