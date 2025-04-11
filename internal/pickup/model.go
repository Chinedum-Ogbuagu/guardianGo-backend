package pickup

type Pickup struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	DropOffID   uint   `json:"drop_off_id"`
	ChurchID    uint   `json:"church_id"`
	PickedUpBy  string `json:"picked_up_by"`
	VerifiedBy  uint   `json:"verified_by"`
	ConfirmedBy string `json:"confirmed_by"`
	PickupTime  string `json:"pickup_time"`
	Notes       string `json:"notes"`
}