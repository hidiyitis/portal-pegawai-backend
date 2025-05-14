package domain

import "time"

type Agenda struct {
	AgendaID     uint      `json:"agenda_id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"not null"  binding:"required"`
	Date         time.Time `json:"date" gorm:"default:CURRENT_TIMESTAMP"  binding:"required"`
	Location     string    `json:"location" gorm:"not null"  binding:"required"`
	Description  string    `json:"description" gorm:"not null"`
	CreatedBy    uint      `json:"created_by" gorm:"not null"`
	Creator      User      `json:"creator" gorm:"foreignKey:CreatedBy;references:NIP"`
	Participants []User    `json:"participants" gorm:"many2many:participants;foreignKey:AgendaID;joinForeignKey:AgendaID;references:NIP;joinReferences:UserNIP"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type InputAgendaRequest struct {
	Title        string `json:"title" binding:"required"`
	Location     string `json:"location" binding:"required"`
	Date         string `json:"date" binding:"required"`
	Description  string `json:"description"`
	Participants []uint `json:"participants"`
	CreatedBy    uint   `json:"-"`
}
