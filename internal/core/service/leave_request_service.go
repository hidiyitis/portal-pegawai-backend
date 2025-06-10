package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/infrastructure/storage"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"github.com/hidiyitis/portal-pegawai/pkg/utils/constants"
)

type LeaveRequestService struct {
	repo        repository.LeaveRequestRepository
	userRepo    repository.UserRepository
	holidayRepo repository.HolidayRepository
	gcpStorage  storage.GCPStorage
}

func NewLeaveRequestService(repo repository.LeaveRequestRepository, userRepo repository.UserRepository, holidayRepo repository.HolidayRepository, gcpStorage storage.GCPStorage) *LeaveRequestService {
	return &LeaveRequestService{repo: repo, userRepo: userRepo, holidayRepo: holidayRepo, gcpStorage: gcpStorage}
}

func (s *LeaveRequestService) CreateLeaveRequest(ctx context.Context, user domain.User, leaveRequest *domain.InputLeaveRequest,
	fileHeader *multipart.FileHeader) (*domain.LeaveRequest, error) {

	startDate, _ := time.Parse(time.RFC3339, leaveRequest.StartDate)
	endDate, _ := time.Parse(time.RFC3339, leaveRequest.EndDate)

	totalHoliday, err := s.holidayRepo.CountHoliday(startDate, endDate)

	if err != nil {
		return nil, err
	}

	totalLeave := calculateDaysBetween(startDate, endDate) - int(totalHoliday)

	if user.LeaveQuota == 0 || user.LeaveQuota < totalLeave {
		return nil, errors.New("user has no leave quota")
	}

	err = validateFile(false, fileHeader)
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
		"portal-pegawai-file",
		file,
		fileHeader.Filename,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	request := &domain.LeaveRequest{
		Title:         leaveRequest.Title,
		StartDate:     startDate,
		ManagerNIP:    leaveRequest.ManagerNIP,
		UserNIP:       user.NIP,
		EndDate:       endDate,
		Description:   leaveRequest.Description,
		FileName:      fileHeader.Filename,
		AttachmentUrl: publicURL,
	}
	result, err := s.repo.CreateLeaveRequest(request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *LeaveRequestService) UpdateLeaveRequest(id uint, user domain.User, leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	// Ambil data existing terlebih dahulu
	existingLeaveRequest, err := s.repo.GetLeaveRequestByID(id)
	if err != nil {
		return nil, fmt.Errorf("leave request not found: %w", err)
	}

	// Validasi status
	isValidStatus := isValidLeaveRequestStatus(leaveRequest.Status)
	if !isValidStatus {
		return nil, fmt.Errorf("invalid leave request status: %s", leaveRequest.Status)
	}

	// Cek permission
	if (leaveRequest.Status == constants.COMPLETED || leaveRequest.Status == constants.REJECTED) &&
		user.NIP != existingLeaveRequest.ManagerNIP {
		return nil, fmt.Errorf("invalid leave request access role status: %s", leaveRequest.Status)
	}

	// Update status
	existingLeaveRequest.Status = leaveRequest.Status

	// LOGIKA ENHANCED: Jika status COMPLETED, hitung dengan precision tinggi
	if leaveRequest.Status == constants.COMPLETED {
		// Hitung working days yang akan dipotong dari kuota
		workingDays, err := s.calculateActualWorkingDays(existingLeaveRequest.StartDate, existingLeaveRequest.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate working days: %w", err)
		}

		userLeave, err := s.userRepo.FindByNIP(existingLeaveRequest.UserNIP)
		if err != nil {
			return nil, err
		}

		// Cek apakah kuota mencukupi
		if userLeave.LeaveQuota < workingDays {
			existingLeaveRequest.Status = constants.REJECTED
			_, err := s.repo.UpdateLeaveRequest(id, existingLeaveRequest)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("rejected: insufficient leave quota. Required: %d days, Available: %d days", workingDays, userLeave.LeaveQuota)
		}

		// Kurangi kuota dengan jumlah working days yang tepat
		userLeave.LeaveQuota -= workingDays
		err = s.userRepo.Update(userLeave)
		if err != nil {
			return nil, err
		}
	}

	// Update leave request
	result, err := s.repo.UpdateLeaveRequest(id, existingLeaveRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (s *LeaveRequestService) calculateActualWorkingDays(startDate time.Time, endDate time.Time) (int, error) {
	if startDate.After(endDate) {
		return 0, fmt.Errorf("start date cannot be after end date")
	}

	// Normalize dates ke midnight untuk consistency
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)

	workingDays := 0
	current := start

	// Ambil daftar holidays dalam range untuk debugging dan verification
	holidays, err := s.holidayRepo.GetHolidaysInRange(start, end)
	if err != nil {
		return 0, fmt.Errorf("failed to get holidays: %w", err)
	}

	// Convert holidays ke map untuk O(1) lookup performance
	holidayMap := make(map[string]bool)
	for _, holiday := range holidays {
		holidayKey := holiday.Date.Format("2006-01-02")
		holidayMap[holidayKey] = true
	}

	// Loop setiap hari dari start sampai end (inclusive)
	for current.Before(end) || current.Equal(end) {
		currentKey := current.Format("2006-01-02")

		// LOGIKA KUNCI: Skip weekend (Saturday = 6, Sunday = 0)
		isWeekend := current.Weekday() == time.Saturday || current.Weekday() == time.Sunday

		// LOGIKA KUNCI: Skip holidays (termasuk yang jatuh pada start/end date)
		isHoliday := holidayMap[currentKey]

		// Hanya hitung sebagai working day jika bukan weekend dan bukan holiday
		if !isWeekend && !isHoliday {
			workingDays++
		}

		// Debug logging untuk transparency (bisa dihapus di production)
		fmt.Printf("Date: %s, Weekend: %t, Holiday: %t, Counted: %t\n",
			currentKey, isWeekend, isHoliday, !isWeekend && !isHoliday)

		// Move to next day
		current = current.AddDate(0, 0, 1)
	}

	return workingDays, nil
}

func (s *LeaveRequestService) GetLeaveRequestDashboard(nip uint) (*domain.LeaveRequestDashboard, error) {
	user, err := s.userRepo.FindByNIP(nip)
	if err != nil {
		return nil, err
	}
	inProgress, err := s.repo.CountLeaveRequestByStatus(nip, constants.IN_PROGRESS)
	if err != nil {
		return nil, err
	}
	rejected, err := s.repo.CountLeaveRequestByStatus(nip, constants.REJECTED)
	if err != nil {
		return nil, err
	}
	cancelled, err := s.repo.CountLeaveRequestByStatus(nip, constants.CANCELLED)
	if err != nil {
		return nil, err
	}
	completed, err := s.repo.CountLeaveRequestByStatus(nip, constants.COMPLETED)
	if err != nil {
		return nil, err
	}

	result := &domain.LeaveRequestDashboard{
		LeaveQuota: user.LeaveQuota,
		InProgress: int(inProgress),
		Cancelled:  int(cancelled),
		Rejected:   int(rejected),
		Completed:  int(completed),
	}

	return result, nil
}

func isValidLeaveRequestStatus(status domain.Status) bool {
	switch status {
	case constants.IN_PROGRESS, constants.REJECTED, constants.CANCELLED, constants.COMPLETED:
		return true
	default:
		return false // It's not valid
	}
}

func calculateDaysBetween(start time.Time, end time.Time) int {
	return int(end.Sub(start).Hours()/24) + 1
}
