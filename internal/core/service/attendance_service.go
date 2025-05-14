package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/infrastructure/storage"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"github.com/hidiyitis/portal-pegawai/pkg/utils/constants"
	"mime/multipart"
)

type AttendanceService struct {
	repo       repository.AttendanceRepository
	gcpStorage storage.GCPStorage
}

func NewAttendanceService(repo repository.AttendanceRepository, gcpStorage storage.GCPStorage) *AttendanceService {
	return &AttendanceService{repo: repo, gcpStorage: gcpStorage}
}

func (s *AttendanceService) CreateAttendance(ctx context.Context, payload domain.InputAttendance, fileHeader *multipart.FileHeader) (*domain.Attendance, error) {
	if !(payload.Status == constants.CLOCK_IN || payload.Status == constants.CLOCK_OUT) {
		return nil, errors.New("invalid attendance status")
	}

	err := validateFile(false, fileHeader)
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	publicURL, err := s.gcpStorage.UploadFile(
		ctx,
		"portal-pegawai-attendance",
		file,
		fileHeader.Filename,
	)
	if err != nil {
		return nil, err
	}
	attendance := &domain.Attendance{
		NIP:       payload.NIP,
		Status:    payload.Status,
		Latitude:  payload.Latitude,
		Longitude: payload.Longitude,
		PhotoUrl:  publicURL,
	}
	if err := s.repo.Create(attendance); err != nil {
		return nil, err
	}
	return attendance, nil
}
