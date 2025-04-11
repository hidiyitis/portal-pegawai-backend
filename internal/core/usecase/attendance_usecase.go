package usecase

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/service"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"mime/multipart"
)

type AttendanceUsecase interface {
	CreateAttendance(attendance domain.InputAttendance, file multipart.File, header *multipart.FileHeader) (*domain.Attendance, error)
}

type attendanceUsecase struct {
	repo    repository.AttendanceRepository
	service *service.AttendanceService
}

func (a attendanceUsecase) CreateAttendance(attendance domain.InputAttendance, file multipart.File, header *multipart.FileHeader) (*domain.Attendance, error) {
	return a.service.CreateAttendance(attendance, file, header)
}

func NewAttendanceUsecase(repo repository.AttendanceRepository, service *service.AttendanceService) AttendanceUsecase {
	return &attendanceUsecase{repo: repo, service: service}
}
