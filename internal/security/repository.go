package security

import "gorm.io/gorm"

type Repository interface {
	Create(db *gorm.DB, flag *SecurityFlag) error
	ListAll(db *gorm.DB) ([]SecurityFlag, error)

}

type repository struct {}

func NewRepository() Repository {
	return &repository{}
}


func (r *repository) Create(db *gorm.DB, flag *SecurityFlag) error {
	return db.Create(flag).Error
}
func (r *repository) ListAll(db *gorm.DB) ([]SecurityFlag, error) {
	var flags []SecurityFlag
	if err := db.Order("created_at desc").Find(&flags).Error; err != nil {
		return nil, err
	}
	return flags, nil
}