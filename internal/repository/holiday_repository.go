package repository

import (
	"time"

	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"gorm.io/gorm"
)

type HolidayRepository interface {
	CountHoliday(startDate time.Time, endDate time.Time) (int64, error)
	GetHolidaysInRange(startDate time.Time, endDate time.Time) ([]domain.Holiday, error) // Tambahan untuk debugging
}

type holidayRepository struct {
	db *gorm.DB
}

func (h holidayRepository) CountHoliday(startDate time.Time, endDate time.Time) (int64, error) {
	// PERBAIKAN UTAMA: Normalize tanggal ke midnight untuk comparison yang akurat
	// Ini seperti memastikan kita membandingkan apel dengan apel, bukan apel dengan jeruk
	normalizedStart := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	normalizedEnd := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)

	var count int64

	// Query yang lebih eksplisit untuk memastikan include start dan end date
	err := h.db.Model(&domain.Holiday{}).
		Where("date >= ? AND date <= ?", normalizedStart, normalizedEnd).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

// TAMBAHAN: Method untuk debugging dan verification
func (h holidayRepository) GetHolidaysInRange(startDate time.Time, endDate time.Time) ([]domain.Holiday, error) {
	normalizedStart := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	normalizedEnd := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)

	var holidays []domain.Holiday
	err := h.db.Where("date >= ? AND date <= ?", normalizedStart, normalizedEnd).
		Order("date ASC").
		Find(&holidays).Error

	return holidays, err
}

func NewHolidayRepository(db *gorm.DB) HolidayRepository {
	return &holidayRepository{db: db}
}
