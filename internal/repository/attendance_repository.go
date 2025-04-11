package repository

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(attendance *domain.Attendance) error
}

type attendanceRepository struct {
	db *gorm.DB
}

func (a attendanceRepository) Create(attendance *domain.Attendance) error {
	return a.db.Create(attendance).Error
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}
