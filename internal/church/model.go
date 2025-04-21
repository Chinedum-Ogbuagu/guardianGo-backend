package church

import "time"

type Church struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	Address	  string    `json:"address"`
	LogoURl      string    `json:"logo_url"`
	CreatedAt    time.Time `json:"created_at"`
}