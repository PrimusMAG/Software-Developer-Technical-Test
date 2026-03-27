package services

import (
	"strings"

	"football-api/internal/domains/teams/dtos"
	teammodels "football-api/internal/domains/teams/models"
	"football-api/internal/domains/teams/repositories"
)

type Service interface {
	Create(req dtos.TeamRequest) (*teammodels.Team, error)
	Update(id uint, req dtos.TeamRequest) (*teammodels.Team, error)
	Delete(id uint) error
	GetByID(id uint) (*teammodels.Team, error)
	List(offset, limit int, search, city string) ([]teammodels.Team, int64, error)
}

type service struct{ repo repositories.Repository }

func New(repo repositories.Repository) Service { return &service{repo: repo} }

func (s *service) Create(req dtos.TeamRequest) (*teammodels.Team, error) {
	team := &teammodels.Team{
		Name:        strings.TrimSpace(req.Name),
		FoundedYear: req.FoundedYear,
		HQAddress:   strings.TrimSpace(req.HQAddress),
		HQCity:      strings.TrimSpace(req.HQCity),
	}
	if err := s.repo.Create(team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *service) Update(id uint, req dtos.TeamRequest) (*teammodels.Team, error) {
	team, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	team.Name = strings.TrimSpace(req.Name)
	team.FoundedYear = req.FoundedYear
	team.HQAddress = strings.TrimSpace(req.HQAddress)
	team.HQCity = strings.TrimSpace(req.HQCity)
	if err := s.repo.Update(team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *service) Delete(id uint) error                      { return s.repo.Delete(id) }
func (s *service) GetByID(id uint) (*teammodels.Team, error) { return s.repo.GetByID(id) }
func (s *service) List(offset, limit int, search, city string) ([]teammodels.Team, int64, error) {
	return s.repo.List(offset, limit, search, city)
}
