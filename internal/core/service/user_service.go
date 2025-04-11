package service

import (
	"errors"
	"fmt"
	"github.com/hidiyitis/portal-pegawai/pkg/utils"
	"strconv"
	"sync"
	"time"

	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo repository.UserRepository
	mu   sync.Mutex
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GenerateNIP(departmentId uint) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	generatedNIP := "100"
	year := time.Now().Year()

	lastData, err := s.repo.FindLast(departmentId)
	if !errors.Is(err, gorm.ErrRecordNotFound) && lastData == nil {
		return "", err
	}
	if lastData == nil || (lastData.CreatedAt.Year() != year) {
		generatedNIP += strconv.Itoa(int(departmentId)) + strconv.Itoa(year)[2:] + "00001"
		return generatedNIP, nil
	}
	println(strconv.FormatUint(uint64(lastData.NIP+1), 10))
	generatedNIP = strconv.FormatUint(uint64(lastData.NIP+1), 10)
	return generatedNIP, nil
}

func (s *UserService) CreateUser(user *domain.User) error {
	fmt.Println(user)
	generatedNIP, err := s.GenerateNIP(user.DepartmentID)
	if err != nil {
		return errors.New("failed to generate ID")
	}
	nip, _ := strconv.Atoi(generatedNIP)
	user.NIP = uint(nip)
	user.Password = strconv.Itoa(nip)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.IsActive = true
	user.Password = string(hashedPassword)
	err = s.repo.Create(user)
	findUser, _ := s.repo.FindByNIP(user.NIP)
	if findUser != nil {
		*user = *findUser
	}
	return err
}

func (s *UserService) LoginUser(user *domain.User) (string, string, string, error) {
	findUser, err := s.repo.FindByNIP(user.NIP)
	if err != nil {
		return "", "", "", errors.New("failed to find user")
	}
	if !findUser.IsActive {
		return "", "", "", errors.New("user is not active")
	}
	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(user.Password))
	if err := bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(user.Password)); err != nil {
		return "", "", "", errors.New("nip and password do not match")
	}
	*user = *findUser
	token, expiredAt, err := utils.GenerateToken(*user)
	if err != nil {
		return "", "", "", errors.New("failed to generate token")
	}
	refreshToken, err := utils.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", "", errors.New("failed to generate refresh token")
	}
	return token, refreshToken, expiredAt, nil
}
