package guardian

type Guardian struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ChurchID    uint   `json:"church_id"`
	Name        string `json:"name"`
	PhoneNumber string `gorm:"unique" json:"phone_number"`
	PhotoURL    string `json:"photo_url"`
	CreatedAt   string `json:"created_at"`
}