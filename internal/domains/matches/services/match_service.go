package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"football-api/internal/domains/matches/dtos"
	matchmodels "football-api/internal/domains/matches/models"
	"football-api/internal/domains/matches/repositories"
	"football-api/internal/shared/constants"
)

type Service interface {
	Create(req dtos.CreateMatchRequest) (*matchmodels.Match, error)
	List(offset, limit int, status string, homeTeamID, awayTeamID uint, dateFrom, dateTo *time.Time) ([]matchmodels.Match, int64, error)
	GetByID(id uint) (*matchmodels.Match, error)
	Delete(id uint) error
	SubmitResult(id uint, req dtos.SubmitResultRequest) (*matchmodels.Match, error)
	RollbackResult(id uint) (*matchmodels.Match, error)
}

type service struct{ repo repositories.Repository }

func New(repo repositories.Repository) Service { return &service{repo: repo} }

func (s *service) Create(req dtos.CreateMatchRequest) (*matchmodels.Match, error) {
	kickoffAt, err := time.Parse(time.RFC3339, req.KickoffAt)
	if err != nil {
		return nil, err
	}
	match := &matchmodels.Match{
		KickoffAt:  kickoffAt,
		HomeTeamID: req.HomeTeamID,
		AwayTeamID: req.AwayTeamID,
		Status:     string(constants.MatchScheduled),
	}
	if err := s.repo.Create(match); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *service) List(offset, limit int, status string, homeTeamID, awayTeamID uint, dateFrom, dateTo *time.Time) ([]matchmodels.Match, int64, error) {
	return s.repo.List(offset, limit, status, homeTeamID, awayTeamID, dateFrom, dateTo)
}

func (s *service) GetByID(id uint) (*matchmodels.Match, error) { return s.repo.GetByID(id) }
func (s *service) Delete(id uint) error                        { return s.repo.Delete(id) }

func (s *service) SubmitResult(id uint, req dtos.SubmitResultRequest) (*matchmodels.Match, error) {
	match, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if match.Status == string(constants.MatchFinished) {
		return nil, errors.New("result already submitted and immutable")
	}

	if req.HomeScore+req.AwayScore != len(req.Goals) {
		return nil, errors.New("total goals must equal homeScore + awayScore")
	}

	goalRows := make([]matchmodels.GoalEvent, 0, len(req.Goals))
	for _, g := range req.Goals {
		if g.TeamID != match.HomeTeamID && g.TeamID != match.AwayTeamID {
			return nil, errors.New("goal team must be either home or away team")
		}
		player, err := s.repo.GetPlayer(g.PlayerID)
		if err != nil {
			return nil, errors.New("invalid player in goal list")
		}
		if player.TeamID != g.TeamID {
			return nil, errors.New("goal player must belong to goal team")
		}
		goalRows = append(goalRows, matchmodels.GoalEvent{
			MatchID:  match.ID,
			TeamID:   g.TeamID,
			PlayerID: g.PlayerID,
			Minute:   g.Minute,
		})
	}

	now := time.Now()
	err = s.repo.WithTx(func(tx *gorm.DB) error {
		match.HomeScore = req.HomeScore
		match.AwayScore = req.AwayScore
		match.Status = string(constants.MatchFinished)
		match.SubmittedAt = &now
		if err := tx.Save(match).Error; err != nil {
			return err
		}
		return s.repo.CreateGoals(tx, goalRows)
	})
	if err != nil {
		return nil, err
	}
	return match, nil
}

func (s *service) RollbackResult(id uint) (*matchmodels.Match, error) {
	match, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if match.Status != string(constants.MatchFinished) {
		return nil, errors.New("match is not finished")
	}
	err = s.repo.WithTx(func(tx *gorm.DB) error {
		if err := s.repo.DeleteGoalsByMatch(tx, match.ID); err != nil {
			return err
		}
		match.Status = string(constants.MatchScheduled)
		match.HomeScore = 0
		match.AwayScore = 0
		match.SubmittedAt = nil
		return tx.Save(match).Error
	})
	if err != nil {
		return nil, err
	}
	return match, nil
}
