package usecase

import (
	"time"

	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/core/service"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
)

type AgendaUsecase interface {
	CreateAgenda(agenda *domain.InputAgendaRequest) (*domain.Agenda, error)
	GetAgendaByID(id uint) (*domain.Agenda, error)
	UpdateAgenda(id uint, agenda *domain.InputAgendaRequest) (*domain.Agenda, error)
	DeleteAgenda(id uint) error
	GetAgendaByDate(nip uint, date time.Time) ([]domain.Agenda, error)
	GetAllAgendas(nip uint) ([]domain.Agenda, error)
}

type agendaUsecase struct {
	repo    repository.AgendaRepository
	service *service.AgendaService
}

func (u agendaUsecase) CreateAgenda(agenda *domain.InputAgendaRequest) (*domain.Agenda, error) {
	return u.service.CreateAgenda(agenda)
}

func (u agendaUsecase) GetAgendaByID(id uint) (*domain.Agenda, error) {
	return u.repo.GetAgendaByID(id)
}

func (u agendaUsecase) UpdateAgenda(id uint, agenda *domain.InputAgendaRequest) (*domain.Agenda, error) {
	return u.service.UpdateAgenda(id, agenda)
}

func (u agendaUsecase) DeleteAgenda(id uint) error {
	return u.repo.DeleteAgenda(id)
}

func (u agendaUsecase) GetAgendaByDate(nip uint, date time.Time) ([]domain.Agenda, error) {
	return u.repo.GetAgendaByDate(nip, date)
}

func NewAgendaUsecase(repo repository.AgendaRepository, service *service.AgendaService) AgendaUsecase {
	return &agendaUsecase{
		repo,
		service,
	}
}

func (u agendaUsecase) GetAllAgendas(nip uint) ([]domain.Agenda, error) {
	return u.repo.GetAllAgendas(nip)
}
