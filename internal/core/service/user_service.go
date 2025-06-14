package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/hidiyitis/portal-pegawai/internal/infrastructure/storage"
	"github.com/hidiyitis/portal-pegawai/pkg/utils"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo       repository.UserRepository
	gcpStorage storage.GCPStorage
	mu         sync.Mutex
}

func NewUserService(repo repository.UserRepository, gcpStorage storage.GCPStorage) *UserService {
	return &UserService{repo: repo, gcpStorage: gcpStorage}
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

func (s *UserService) UploadAvatar(ctx context.Context, user *domain.User, fileHeader *multipart.FileHeader) (*domain.User, error) {
	if fileHeader == nil {
		return nil, fmt.Errorf("fileHeader is nil")
	}
	user, err := s.repo.FindByNIP(user.NIP)
	if err != nil {
		return nil, err
	}
	if err := validateFile(true, fileHeader); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
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
		"portal-pegawai-images",
		file,
		fileHeader.Filename,
	)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	user.PhotoUrl = publicURL

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdatePasswordUser(user domain.User, payload domain.UpdateUserPassword) (*domain.User, error) {
	findUser, err := s.repo.FindByNIP(user.NIP)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(payload.CurrentPassword))
	if err != nil {
		return nil, errors.New("password does not match")
	}
	if payload.NewPassword != payload.ConfirmPassword {
		return nil, errors.New("new password does not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	findUser.Password = string(hashedPassword)

	err = s.repo.Update(findUser)
	if err != nil {
		return nil, err
	}
	return findUser, nil
}

func validateFile(isImage bool, fileHeader *multipart.FileHeader) error {
	const maxFileSize = 5 << 20 // 5MB
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file too large, max size is %dMB", maxFileSize>>20)
	}

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	ext := filepath.Ext(fileHeader.Filename)
	if !allowedExts[ext] && isImage {
		return fmt.Errorf("unsupported file type, allowed: jpg, jpeg, png")
	}

	return nil
}
