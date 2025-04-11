package usecase

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/service"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	GetUsers() ([]domain.User, error)
	GetUserByNIP(nip uint) (*domain.User, error)
	UpdateUser(user *domain.User) error
	LoginUser(user *domain.User) (string, string, string, error)
}

type userUsecase struct {
	repo    repository.UserRepository
	service *service.UserService
}

// CreateUser implements UserUsecase.
func (u *userUsecase) CreateUser(user *domain.User) error {
	return u.service.CreateUser(user)
}

// GetUserByID implements UserUsecase.
func (u *userUsecase) GetUserByNIP(nip uint) (*domain.User, error) {
	return u.repo.FindByNIP(nip)
}

// GetUsers implements UserUsecase.
func (u *userUsecase) GetUsers() ([]domain.User, error) {
	return u.repo.FindAll()
}

// UpdateUser implements UserUsecase.
func (u *userUsecase) UpdateUser(user *domain.User) error {
	return u.repo.Update(user)
}

// LoginUser implements UserUsecase
func (u *userUsecase) LoginUser(user *domain.User) (string, string, string, error) {
	return u.service.LoginUser(user)
}

func NewUserUsecase(repo repository.UserRepository, service *service.UserService) UserUsecase {
	return &userUsecase{repo: repo, service: service}
}
