package domain

type Participant struct {
	UserNIP  uint `gorm:"primaryKey"`
	AgendaID uint `gorm:"primaryKey"`
}
