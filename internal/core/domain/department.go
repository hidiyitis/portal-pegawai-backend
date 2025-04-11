package domain

import "time"

type Department struct {
	DepartmentId uint      `json:"department_id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
