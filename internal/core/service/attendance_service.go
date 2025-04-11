package service

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"github.com/hidiyitis/portal-pegawai/pkg/utils"
	"mime/multipart"
)

type AttendanceService struct {
	repo repository.AttendanceRepository
}

func NewAttendanceService(repo repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

func (s *AttendanceService) CreateAttendance(payload domain.InputAttendance, file multipart.File, header *multipart.FileHeader) (*domain.Attendance, error) {
	filename, err := utils.SaveFile(file, header)
	if err != nil {
		return nil, err
	}
	attendance := &domain.Attendance{
		AttendanceID: payload.NIP,
		Status:       payload.Status,
		PhotoUrl:     filename,
	}
	if err := s.repo.Create(attendance); err != nil {
		return nil, err
	}
	return attendance, nil
}
