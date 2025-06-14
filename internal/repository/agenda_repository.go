package repository

import (
	"fmt"
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"gorm.io/gorm"
	"time"
)

type AgendaRepository interface {
	CreateAgenda(agenda *domain.Agenda, participants []uint) (*domain.Agenda, error)
	GetAgendaByID(id uint) (*domain.Agenda, error)
	UpdateAgenda(agenda *domain.Agenda, participants []uint) (*domain.Agenda, error)
	GetAgendaByDate(nip uint, date time.Time) ([]domain.Agenda, error)
	DeleteAgenda(id uint) error
}

type agendaRepository struct {
	db *gorm.DB
}

func (a agendaRepository) CreateAgenda(agenda *domain.Agenda, participants []uint) (*domain.Agenda, error) {
	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(agenda).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	var users []domain.User
	if err := tx.Where("NIP IN ?", participants).Find(&users).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Model(agenda).Association("Participants").Append(users); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := a.db.Preload("Creator").First(agenda, agenda.AgendaID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload agenda: %w", err)
	}
	return agenda, nil
}

func (a agendaRepository) GetAgendaByID(id uint) (*domain.Agenda, error) {
	agenda := &domain.Agenda{}
	err := a.db.Preload("Participants").First(&agenda, id).Error
	if err != nil {
		return nil, err
	}
	return agenda, nil

}

func (a agendaRepository) UpdateAgenda(agenda *domain.Agenda, participants []uint) (*domain.Agenda, error) {
	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Save(agenda).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	var users []domain.User
	if err := tx.Where("NIP IN ?", participants).Find(&users).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Model(agenda).Association("Participants").Replace(users); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return agenda, nil
}

func (a agendaRepository) GetAgendaByDate(nip uint, date time.Time) ([]domain.Agenda, error) {
	var agenda []domain.Agenda
	print(nip)
	err := a.db.Preload("Participants").Preload("Creator").Joins("JOIN participants ON participants.agenda_id = agendas.agenda_id").
		Where("DATE(date) = ? AND participants.user_n_ip = ?  AND date > ?", date.Format(time.DateOnly), nip, time.Now().UTC()).
		Order("date ASC").Find(&agenda).Error
	if err != nil {
		return nil, err
	}
	return agenda, nil
}

func (a agendaRepository) DeleteAgenda(id uint) error {
	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("agenda_id = ?", id).Delete(&domain.Participant{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("agenda_id = ?", id).Delete(&domain.Agenda{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func NewAgendaRepository(db *gorm.DB) AgendaRepository {
	return &agendaRepository{
		db: db,
	}
}
