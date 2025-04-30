package church

import (
	"time"

	"github.com/gofrs/uuid"
)

type Church struct {
	ID        	 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string    `json:"name"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	Address	     string    `json:"address"`
	LogoURl      string    `json:"logo_url"`
	CreatedAt    time.Time `json:"created_at"`
}