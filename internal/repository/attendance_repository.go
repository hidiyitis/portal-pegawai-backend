package repository

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(attendance *domain.Attendance) error
	GetLastAttendance(nip uint) (*domain.Attendance, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func (a attendanceRepository) GetLastAttendance(nip uint) (*domain.Attendance, error) {
	attendance := &domain.Attendance{}
	err := a.db.Last(&attendance, "n_ip = ?", nip).Error
	if err != nil {
		return nil, err
	}
	return attendance, nil
}

func (a attendanceRepository) Create(attendance *domain.Attendance) error {
	return a.db.Create(attendance).Error
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}
