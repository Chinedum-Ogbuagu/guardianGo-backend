package dropoff

import (
	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(100000 + rand.Intn(899999)) 
}

type Service interface {
	Create(db *gorm.DB, d *DropOff) error 
	GetByID(db *gorm.DB, id uint ) (*DropOff, error)
}

type service struct {
	repo Repository
}

func NewService (r Repository) Service {
	return &service{repo: r}
}
func (s *service) Create(db *gorm.DB, d *DropOff) error {
	d.UniqueCode = generateCode()
	d.DropOffTime = time.Now().Format(time.RFC3339)
	return s.repo.Create(db, d)
}

func (s *service) GetByID(db *gorm.DB, id uint ) (*DropOff, error) {
	return s.repo.GetByID(db, id)
}