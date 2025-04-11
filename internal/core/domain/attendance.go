package domain

import (
	"time"
)

type AttendenceStatus string

type Attendance struct {
	AttendanceID uint             `json:"attendance_id" gorm:"primary_key;AUTO_INCREMENT"`
	NIP          uint             `json:"nip" gorm:"foreignKey"`
	PhotoUrl     string           `json:"photo_url" gorm:"column:photo_url" binding:"required"`
	Status       AttendenceStatus `json:"status" gorm:"column:status" binding:"required"`
	CreatedAt    time.Time        `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type InputAttendance struct {
	NIP    uint
	Status AttendenceStatus `form:"status" binding:"required"`
}
