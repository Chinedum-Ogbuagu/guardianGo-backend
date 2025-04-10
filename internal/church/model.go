package church

type Church struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `json:"name"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
	LogoURl      string `json:"logo_url"`
	CreatedAt    string `json:"created_at"`
}