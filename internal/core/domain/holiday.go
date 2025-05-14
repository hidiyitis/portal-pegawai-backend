package domain

import "time"

type Holiday struct {
	HolidayID   uint      `json:"holiday_id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"  binding:"required"`
	Date        time.Time `json:"date" gorm:"not null"  binding:"required"`
	Description string    `json:"description" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
