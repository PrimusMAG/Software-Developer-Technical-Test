package services

import (
	"strings"

	"football-api/internal/domains/players/dtos"
	playermodels "football-api/internal/domains/players/models"
	"football-api/internal/domains/players/repositories"
)

type Service interface {
	Create(req dtos.PlayerRequest) (*playermodels.Player, error)
	Update(id uint, req dtos.PlayerRequest) (*playermodels.Player, error)
	Delete(id uint) error
	GetByID(id uint) (*playermodels.Player, error)
	List(offset, limit int, search, position string, teamID uint) ([]playermodels.Player, int64, error)
}

type service struct{ repo repositories.Repository }

func New(repo repositories.Repository) Service { return &service{repo: repo} }

func (s *service) Create(req dtos.PlayerRequest) (*playermodels.Player, error) {
	player := &playermodels.Player{
		TeamID:       req.TeamID,
		Name:         strings.TrimSpace(req.Name),
		HeightCM:     req.HeightCM,
		WeightKG:     req.WeightKG,
		Position:     req.Position,
		JerseyNumber: req.JerseyNumber,
	}
	if err := s.repo.Create(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (s *service) Update(id uint, req dtos.PlayerRequest) (*playermodels.Player, error) {
	player, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	player.TeamID = req.TeamID
	player.Name = strings.TrimSpace(req.Name)
	player.HeightCM = req.HeightCM
	player.WeightKG = req.WeightKG
	player.Position = req.Position
	player.JerseyNumber = req.JerseyNumber
	if err := s.repo.Update(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (s *service) Delete(id uint) error                          { return s.repo.Delete(id) }
func (s *service) GetByID(id uint) (*playermodels.Player, error) { return s.repo.GetByID(id) }
func (s *service) List(offset, limit int, search, position string, teamID uint) ([]playermodels.Player, int64, error) {
	return s.repo.List(offset, limit, search, position, teamID)
}
