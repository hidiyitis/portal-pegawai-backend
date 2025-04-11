package repository

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindAll() ([]domain.User, error)
	FindByNIP(nip uint) (*domain.User, error)
	Update(user *domain.User) error
	FindLast(departmentId uint) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func (u userRepository) Create(user *domain.User) error {
	return u.db.Create(&user).Error
}

func (u userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := u.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u userRepository) FindByNIP(nip uint) (*domain.User, error) {
	user := &domain.User{}
	err := u.db.Preload("Department").Where("nip=?", nip).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u userRepository) Update(user *domain.User) error {
	return u.db.Save(user).Error
}

func (u userRepository) FindLast(departmentId uint) (*domain.User, error) {
	var user domain.User
	err := u.db.Last(&user, "department_id = ?", departmentId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
