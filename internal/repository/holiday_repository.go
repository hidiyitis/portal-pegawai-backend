package repository

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"gorm.io/gorm"
	"time"
)

type HolidayRepository interface {
	CountHoliday(startDate time.Time, endDate time.Time) (int64, error)
}

type holidayRepository struct {
	db *gorm.DB
}

func (h holidayRepository) CountHoliday(startDate time.Time, endDate time.Time) (int64, error) {
	var count int64
	err := h.db.Model(&domain.Holiday{}).Where("date BETWEEN ? AND ?", startDate, endDate).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func NewHolidayRepository(db *gorm.DB) HolidayRepository {
	return &holidayRepository{db: db}
}
