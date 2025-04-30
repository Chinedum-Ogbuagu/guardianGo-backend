package user

import (
	"time"

	uuid "github.com/gofrs/uuid"
)

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleChurchAdmin Role = "church_admin"
	RoleAttendant  Role = "attendant"
	RoleSecurity Role = "security"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name       string    `json:"name"`
	Phone      string    `gorm:"uniqueIndex" json:"phone"`
	Role       Role      `json:"role"`
	ChurchID  *uuid.UUID `json:"church_id"`
	CreatedAt  time.Time `json:"created_at"`
}
