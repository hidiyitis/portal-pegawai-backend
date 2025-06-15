package usecase

import (
	"context"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/service"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"mime/multipart"
)

type AttendanceUsecase interface {
	CreateAttendance(ctx context.Context, attendance domain.InputAttendance, fileHeader *multipart.FileHeader) (*domain.Attendance, error)
	GetLastAttendance(nip uint) (*domain.Attendance, error)
}

type attendanceUsecase struct {
	repo    repository.AttendanceRepository
	service *service.AttendanceService
}

func (a attendanceUsecase) GetLastAttendance(nip uint) (*domain.Attendance, error) {
	return a.repo.GetLastAttendance(nip)
}

func (a attendanceUsecase) CreateAttendance(ctx context.Context, attendance domain.InputAttendance, fileHeader *multipart.FileHeader) (*domain.Attendance, error) {
	return a.service.CreateAttendance(ctx, attendance, fileHeader)
}

func NewAttendanceUsecase(repo repository.AttendanceRepository, service *service.AttendanceService) AttendanceUsecase {
	return &attendanceUsecase{repo: repo, service: service}
}
