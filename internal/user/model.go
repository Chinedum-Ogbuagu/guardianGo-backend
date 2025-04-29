package user

import "time"

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleAttendant  Role = "attendant"
)

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	Phone      string    `gorm:"uniqueIndex" json:"phone"`
	Role       Role      `json:"role"`
	ChurchID   uint      `json:"church_id"`
	CreatedAt  time.Time `json:"created_at"`
}
