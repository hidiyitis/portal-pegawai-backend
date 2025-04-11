package usecase

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/service"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"mime/multipart"
)

type LeaveRequestUsecase interface {
	CreateLeaveRequest(leaveRequest *domain.InputLeaveRequest, file multipart.File, header *multipart.FileHeader) (*domain.LeaveRequest, error)
}

type leaveRequestUsercase struct {
	repo    repository.LeaveRequestRepository
	service *service.LeaveRequestService
}

func (u leaveRequestUsercase) CreateLeaveRequest(leaveRequest *domain.InputLeaveRequest, file multipart.File, header *multipart.FileHeader) (*domain.LeaveRequest, error) {
	return u.service.CreateLeaveRequest(leaveRequest, file, header)
}

func NewLeaveRequestUsecase(repo repository.LeaveRequestRepository, service *service.LeaveRequestService) LeaveRequestUsecase {
	return &leaveRequestUsercase{
		repo:    repo,
		service: service,
	}
}
