package domain

import (
	"time"
)

type Status string

type LeaveRequest struct {
	LeaveID       uint      `json:"leave_id" gorm:"primary_key"`
	Title         string    `json:"title" gorm:"size:255;not null"`
	UserNIP       uint      `json:"user_nip" gorm:"column:user_nip; not null"`
	ManagerNIP    uint      `json:"manager_nip" gorm:"column:manager_nip; not null"`
	StartDate     time.Time `json:"start_date" gorm:"column:start_date; not null"`
	EndDate       time.Time `json:"end_date" gorm:"column:end_date; not null"`
	FileName      string    `json:"file_name" gorm:"column:file_name; not null"`
	AttachmentUrl string    `json:"attachment_url" gorm:"column:attachment_url;"`
	Description   string    `json:"description" gorm:"column:description;"`
	Status        Status    `json:"status" gorm:"column:status; default:'IN_PROGRESS';"`
}

type InputLeaveRequest struct {
	Title       string `form:"title" json:"title" binding:"required"`
	StartDate   string `form:"start_date" json:"start_date" binding:"required"`
	EndDate     string `form:"end_date" json:"end_date" binding:"required"`
	ManagerNIP  uint   `form:"manager_nip" json:"manager_nip" binding:"required" `
	Description string `form:"description" json:"description"`
	UserNIP     uint
}

type LeaveRequestDashboard struct {
	LeaveQuota int `json:"leave_quota"`
	InProgress int `json:"in_progress"`
	Cancelled  int `json:"cancelled"`
	Rejected   int `json:"rejected"`
	Completed  int `json:"completed"`
}
