package domain

import "time"

type User struct {
	NIP          uint       `json:"nip" gorm:"primaryKey; column:nip"`
	Name         string     `json:"name" gorm:"not null"`
	Password     string     `json:"-" gorm:"not null"`
	LeaveQuota   int        `json:"leave_quota" gorm:"not null default 0"`
	PhotoUrl     string     `json:"photo_url" gorm:"not null"`
	Role         string     `json:"role" gorm:"not null"`
	DepartmentID uint       `json:"department_id" gorm:"not null"`
	Department   Department `json:"department" gorm:"foreignKey:DepartmentID; references department(id)"`
	IsActive     bool       `json:"is_active" gorm:"not null default true"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type UpdateUserPassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
