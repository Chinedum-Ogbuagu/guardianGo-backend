package guardian

import "gorm.io/gorm"

type Repository interface {
	FindByPhone(db *gorm.DB, phone string) (*Guardian, error)
	Create(db *gorm.DB, g *Guardian) error
	GetChildrenByGuardianPhone(db *gorm.DB, phone string) ([]ChildInfo, error)
}
type ChildInfo struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Note  string `json:"note"`
	Bag   bool   `json:"bag"`
}
type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindByPhone(db *gorm.DB, phone string) (*Guardian, error) {
	var guardian Guardian
	if err := db.Where("phone = ?", phone).First(&guardian).Error; err != nil {
		return nil, err
	}
	return &guardian, nil
}

func (r *repository) Create(db *gorm.DB, g *Guardian) error {
	return db.Create(g).Error
}
func (r *repository) GetChildrenByGuardianPhone(db *gorm.DB, phone string) ([]ChildInfo, error) {
	var results []ChildInfo

	query := `
		SELECT c.id, c.name, d.class, d.note, d.bag_status
		FROM children c
		JOIN guardians g ON g.id = c.guardian_id
		LEFT JOIN LATERAL (
			SELECT class, note, bag_status
			FROM drop_offs
			WHERE drop_offs.child_id = c.id
			ORDER BY drop_off_time DESC
			LIMIT 1
		) d ON true
		WHERE g.phone = ?
	`

	err := db.Raw(query, phone).Scan(&results).Error
	return results, err
}

