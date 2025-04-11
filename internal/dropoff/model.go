package dropoff

type DropOff struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ChildID     uint   `json:"child_id"`
	ParentID    uint   `json:"parent_id"`
	ChurchID    uint   `json:"church_id"`
	UniqueCode  string `gorm:"unique" json:"unique_code"`
	ChildAge    int    `json:"child_age"`
	ChildClass  string `json:"child_class"`
	BagStatus   bool   `json:"bag_status"`
	Note        string `json:"note"`
	DropOffTime string `json:"drop_off_time"`
}
