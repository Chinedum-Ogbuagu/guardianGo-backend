package user

type User struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ChurchID    uint   `json:"church_id"`
	Name        string `json:"name"`
	PhoneNumber string `gorm:"unique" json:"phone_number"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
}
