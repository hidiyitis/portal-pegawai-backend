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
	"time"
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
	isValidStatus := isValidLeaveRequestStatus(leaveRequest.Status)
	if !isValidStatus {
		return nil, fmt.Errorf("invalid leave request status: %s", leaveRequest.Status)
	}
	print(user.NIP)

	if (leaveRequest.Status == constants.COMPLETED || leaveRequest.Status == constants.REJECTED) && user.NIP != leaveRequest.ManagerNIP {
		return nil, fmt.Errorf("invalid leave request access role status: %s", leaveRequest.Status)
	}

	if leaveRequest.Status == constants.COMPLETED {
		totalHoliday, err := s.holidayRepo.CountHoliday(leaveRequest.StartDate, leaveRequest.EndDate)
		if err != nil {
			return nil, err
		}

		userLeave, err := s.userRepo.FindByNIP(leaveRequest.UserNIP)
		if err != nil {
			return nil, err
		}

		totalLeave := calculateDaysBetween(leaveRequest.StartDate, leaveRequest.EndDate) - int(totalHoliday)
		if userLeave.LeaveQuota < totalLeave {
			leaveRequest.Status = constants.REJECTED
			_, err := s.repo.UpdateLeaveRequest(id, leaveRequest)
			if err != nil {
				return nil, err
			}

			return nil, errors.New("rejected user has no leave quota")
		}

		userLeave.LeaveQuota -= totalLeave
		err = s.userRepo.Update(userLeave)
		if err != nil {
			return nil, err
		}
	}

	result, err := s.repo.UpdateLeaveRequest(id, leaveRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
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
