package repository

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"gorm.io/gorm"
)

type LeaveRequestRepository interface {
	CreateLeaveRequest(leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error)
	GetLeaveRequestByID(id uint) (*domain.LeaveRequest, error)
	UpdateLeaveRequest(id uint, leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error)
	GetLeaveRequests(nip uint) ([]*domain.LeaveRequest, error)
	GetLeaveApprovals(managerNIP uint) ([]*domain.LeaveRequest, error)
}

type leaveRequestRepository struct {
	db *gorm.DB
}

func (l leaveRequestRepository) CreateLeaveRequest(leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	err := l.db.Create(leaveRequest).Error
	if err != nil {
		return nil, err
	}
	return leaveRequest, nil
}

func (l leaveRequestRepository) GetLeaveRequestByID(id uint) (*domain.LeaveRequest, error) {
	leaveRequest := &domain.LeaveRequest{}
	if err := l.db.Where("id = ?", id).First(leaveRequest).Error; err != nil {
		return nil, err
	}
	return leaveRequest, nil
}

func (l leaveRequestRepository) UpdateLeaveRequest(id uint, leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	if err := l.db.Where("id = ?", id).Updates(leaveRequest).Error; err != nil {
		return nil, err
	}
	return leaveRequest, nil
}

func (l leaveRequestRepository) GetLeaveRequests(nip uint) ([]*domain.LeaveRequest, error) {
	var leaveRequests []*domain.LeaveRequest
	l.db.Where("user_nip = ?", nip).Find(&leaveRequests)
	return leaveRequests, nil
}

func (l leaveRequestRepository) GetLeaveApprovals(managerNIP uint) ([]*domain.LeaveRequest, error) {
	var leaveRequests []*domain.LeaveRequest
	l.db.Where("manager_nip = ?", managerNIP).Find(&leaveRequests).Order("leave_requests.id asc")
	return leaveRequests, nil
}

func NewLeaveRequestRepository(db *gorm.DB) LeaveRequestRepository {
	return &leaveRequestRepository{
		db: db,
	}
}
