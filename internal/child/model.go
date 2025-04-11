package child

type Child struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ParentID  uint   `json:"parent_id"`
	ChurchID  uint   `json:"church_id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	BagStatus bool   `json:"bag_status"`
	PhotoURL  string `json:"photo_url"`
	CreatedAt string `json:"created_at"`
}