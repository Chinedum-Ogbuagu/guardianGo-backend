package security

type SecurityFlag struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	DropOffID uint   `json:"drop_off_id"`
	FlaggedBy string `json:"flagged_by"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
}
