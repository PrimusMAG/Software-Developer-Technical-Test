package services

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"football-api/internal/domains/matches/dtos"
	matchmodels "football-api/internal/domains/matches/models"
	playermodels "football-api/internal/domains/players/models"
	"football-api/internal/shared/constants"
)

type mockRepo struct {
	match *matchmodels.Match
}

func (m *mockRepo) Create(match *matchmodels.Match) error { return nil }
func (m *mockRepo) Update(match *matchmodels.Match) error { return nil }
func (m *mockRepo) Delete(id uint) error                  { return nil }
func (m *mockRepo) GetByID(id uint) (*matchmodels.Match, error) {
	if m.match == nil {
		return nil, errors.New("not found")
	}
	return m.match, nil
}
func (m *mockRepo) List(offset, limit int, status string, homeTeamID, awayTeamID uint, dateFrom, dateTo *time.Time) ([]matchmodels.Match, int64, error) {
	return nil, 0, nil
}
func (m *mockRepo) GetPlayer(id uint) (*playermodels.Player, error) {
	return &playermodels.Player{ID: id, TeamID: 1}, nil
}
func (m *mockRepo) CreateGoals(tx *gorm.DB, goals []matchmodels.GoalEvent) error { return nil }
func (m *mockRepo) DeleteGoalsByMatch(tx *gorm.DB, matchID uint) error           { return nil }
func (m *mockRepo) WithTx(fn func(tx *gorm.DB) error) error                      { return fn(&gorm.DB{}) }

func TestSubmitResult_Immutable(t *testing.T) {
	svc := New(&mockRepo{
		match: &matchmodels.Match{
			ID:         1,
			HomeTeamID: 1,
			AwayTeamID: 2,
			Status:     string(constants.MatchFinished),
		},
	})
	_, err := svc.SubmitResult(1, dtos.SubmitResultRequest{HomeScore: 1, AwayScore: 0, Goals: []dtos.GoalInput{{TeamID: 1, PlayerID: 1, Minute: 10}}})
	if err == nil {
		t.Fatal("expected immutable error, got nil")
	}
}

func TestSubmitResult_ScoreGoalMismatch(t *testing.T) {
	svc := New(&mockRepo{
		match: &matchmodels.Match{
			ID:         1,
			HomeTeamID: 1,
			AwayTeamID: 2,
			Status:     string(constants.MatchScheduled),
		},
	})
	_, err := svc.SubmitResult(1, dtos.SubmitResultRequest{HomeScore: 2, AwayScore: 0, Goals: []dtos.GoalInput{{TeamID: 1, PlayerID: 1, Minute: 10}}})
	if err == nil {
		t.Fatal("expected mismatch error, got nil")
	}
}
