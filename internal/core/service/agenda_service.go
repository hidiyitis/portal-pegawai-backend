package service

import (
	"github.com/hidiyitis/portal-pegawai/internal/core/domain"
	"github.com/hidiyitis/portal-pegawai/internal/repository"
	"time"
)

type AgendaService struct {
	repoAgenda repository.AgendaRepository
	repoUser   repository.UserRepository
}

func NewAgendaService(repoAgenda repository.AgendaRepository, repoUser repository.UserRepository) *AgendaService {
	return &AgendaService{repoAgenda: repoAgenda, repoUser: repoUser}
}

func (s *AgendaService) CreateAgenda(payload *domain.InputAgendaRequest) (*domain.Agenda, error) {
	date, err := time.Parse(time.RFC3339, payload.Date)
	if err != nil {
		return nil, err
	}
	agenda := &domain.Agenda{
		Title:       payload.Title,
		Description: payload.Description,
		Date:        date,
		Location:    payload.Location,
		CreatedBy:   payload.CreatedBy,
	}

	result, err := s.repoAgenda.CreateAgenda(agenda, payload.Participants)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AgendaService) UpdateAgenda(id uint, payload *domain.InputAgendaRequest) (*domain.Agenda, error) {
	findAgenda, err := s.repoAgenda.GetAgendaByID(id)
	if err != nil {
		return nil, err
	}
	if payload.Title != "" {
		findAgenda.Title = payload.Title
	}
	if payload.Description != "" {
		findAgenda.Description = payload.Description
	}
	if payload.Location != "" {
		findAgenda.Location = payload.Location
	}
	if payload.Date != "" {
		date, err := time.Parse(time.RFC3339, payload.Date)
		if err != nil {
			return nil, err
		}
		findAgenda.Date = date
	}
	result, err := s.repoAgenda.UpdateAgenda(findAgenda, payload.Participants)
	if err != nil {
		return nil, err
	}
	return result, nil
}
