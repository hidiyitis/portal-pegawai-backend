package repository

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"gorm.io/gorm"
)

type StatusCount struct {
	Status domain.Status
	Count  int64
}

type LeaveRequestRepository interface {
	CreateLeaveRequest(leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error)
	GetLeaveRequestByID(id uint) (*domain.LeaveRequest, error)
	UpdateLeaveRequest(id uint, leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error)
	GetLeaveRequests(nip uint, status string) ([]*domain.LeaveRequest, error)
	GetLeaveApprovals(managerNIP uint) ([]*domain.LeaveRequest, error)
	CountLeaveRequestByStatus(nip uint, status domain.Status) (int64, error)
}

type leaveRequestRepository struct {
	db *gorm.DB
}

func (l leaveRequestRepository) CountLeaveRequestByStatus(nip uint, status domain.Status) (int64, error) {
	var count int64
	err := l.db.Model(&domain.LeaveRequest{}).Where("user_nip = ? AND status = ?", nip, status).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
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
	if err := l.db.Where("leave_id = ?", id).First(leaveRequest).Error; err != nil {
		return nil, err
	}
	return leaveRequest, nil
}

func (l leaveRequestRepository) UpdateLeaveRequest(id uint, leaveRequest *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	if err := l.db.Where("leave_id = ?", id).Updates(leaveRequest).Error; err != nil {
		return nil, err
	}
	return leaveRequest, nil
}

func (l leaveRequestRepository) GetLeaveRequests(nip uint, status string) ([]*domain.LeaveRequest, error) {
	var leaveRequests []*domain.LeaveRequest
	if status != "" {
		l.db.Where("user_nip = ? AND status = ?", nip, status).Find(&leaveRequests)
		return leaveRequests, nil
	}
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
