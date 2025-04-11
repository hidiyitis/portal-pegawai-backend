package service

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"github.com/hidiyitis/portal-pegawai/pkg/utils"
	"mime/multipart"
	"time"
)

type LeaveRequestService struct {
	repo repository.LeaveRequestRepository
}

func NewLeaveRequestService(repo repository.LeaveRequestRepository) *LeaveRequestService {
	return &LeaveRequestService{repo: repo}
}

func (s *LeaveRequestService) CreateLeaveRequest(leaveRequest *domain.InputLeaveRequest,
	file multipart.File, header *multipart.FileHeader) (*domain.LeaveRequest, error) {
	filename, err := utils.SaveFile(file, header)
	if err != nil {
		return nil, err
	}
	date, _ := time.Parse(time.RFC3339, leaveRequest.Date)
	request := &domain.LeaveRequest{
		Title:         leaveRequest.Title,
		Date:          date,
		Description:   leaveRequest.Description,
		AttachmentUrl: filename,
	}
	result, err := s.repo.CreateLeaveRequest(request)
	if err != nil {
		return nil, err
	}
	return result, nil
}
