package domain

import (
	"time"
)

type AttendenceStatus string

type Attendance struct {
	AttendanceID uint             `json:"attendance_id" gorm:"primary_key;AUTO_INCREMENT"`
	NIP          uint             `json:"nip"`
	PhotoUrl     string           `json:"photo_url" gorm:"column:photo_url" binding:"required"`
	Status       AttendenceStatus `json:"status" gorm:"column:status" binding:"required"`
	Latitude     float64          `json:"latitude" gorm:"type:decimal(10,8)" binding:"required"`
	Longitude    float64          `json:"longitude" gorm:"type:decimal(11,8)" binding:"required"`
	CreatedAt    time.Time        `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type InputAttendance struct {
	NIP       uint
	Status    AttendenceStatus `form:"status" binding:"required"`
	Latitude  float64          `form:"latitude" binding:"required"`
	Longitude float64          `form:"longitude" binding:"required"`
}
