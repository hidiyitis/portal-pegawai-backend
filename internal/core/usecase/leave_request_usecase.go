package usecase

import (
	"context"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/service"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"mime/multipart"
)

type LeaveRequestUsecase interface {
	CreateLeaveRequest(ctx context.Context, user domain.User, leaveRequest *domain.InputLeaveRequest, fileHeader *multipart.FileHeader) (*domain.LeaveRequest, error)
	UpdateLeaveRequest(id uint, user domain.User, leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error)
	GetLeaveRequests(nip uint, status string) ([]*domain.LeaveRequest, error)
	GetLeaveRequestDashboard(nip uint) (*domain.LeaveRequestDashboard, error)
}

type leaveRequestUsercase struct {
	repo    repository.LeaveRequestRepository
	service *service.LeaveRequestService
}

func (u leaveRequestUsercase) GetLeaveRequestDashboard(nip uint) (*domain.LeaveRequestDashboard, error) {
	return u.service.GetLeaveRequestDashboard(nip)
}

func (u leaveRequestUsercase) GetLeaveRequests(nip uint, status string) ([]*domain.LeaveRequest, error) {
	return u.repo.GetLeaveRequests(nip, status)
}

func (u leaveRequestUsercase) UpdateLeaveRequest(id uint, user domain.User, leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	return u.service.UpdateLeaveRequest(id, user, leaveRequest)
}

func (u leaveRequestUsercase) CreateLeaveRequest(ctx context.Context, user domain.User, leaveRequest *domain.InputLeaveRequest, fileHeader *multipart.FileHeader) (*domain.LeaveRequest, error) {
	return u.service.CreateLeaveRequest(ctx, user, leaveRequest, fileHeader)
}

func NewLeaveRequestUsecase(repo repository.LeaveRequestRepository, service *service.LeaveRequestService) LeaveRequestUsecase {
	return &leaveRequestUsercase{
		repo:    repo,
		service: service,
	}
}
