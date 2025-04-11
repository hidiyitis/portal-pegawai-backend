package domain

import (
	"time"
)

type Status string

type LeaveRequest struct {
	LeaveID       uint      `gorm:"primary_key"`
	Title         string    `json:"title" gorm:"size:255;not null"`
	UserNIP       uint      `json:"user_nip" gorm:"column:user_nip; not null"`
	ManagerNIP    uint      `json:"manager_nip" gorm:"column:manager_nip; not null"`
	Date          time.Time `json:"date" gorm:"column:date; not null"`
	AttachmentUrl string    `json:"attachment_url" gorm:"column:attachment_url;"`
	Description   string    `json:"description" gorm:"column:description;"`
	Status        Status    `json:"status" gorm:"column:status; default:'IN_PROGRESS';"`
}

type InputLeaveRequest struct {
	Title       string `form:"title" json:"title" binding:"required"`
	Date        string `form:"date" json:"date" binding:"required"`
	ManagerNIP  uint   `form:"manager_nip" json:"manager_nip" binding:"required" `
	Description string `form:"description" json:"description" binding:"optional"`
}
