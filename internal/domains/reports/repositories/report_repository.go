package repositories

import (
	"gorm.io/gorm"

	matchmodels "football-api/internal/domains/matches/models"
	teammodels "football-api/internal/domains/teams/models"
)

type GoalRow struct {
	MatchID    uint
	PlayerName string
}

type Repository interface {
	ListFinishedMatches() ([]matchmodels.Match, error)
	GetTeamsByIDs(ids []uint) (map[uint]string, error)
	ListGoalsByMatch(matchID uint) ([]GoalRow, error)
}

type repository struct{ db *gorm.DB }

func New(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) ListFinishedMatches() ([]matchmodels.Match, error) {
	var matches []matchmodels.Match
	err := r.db.Where("status = ?", "finished").Order("kickoff_at ASC, id ASC").Find(&matches).Error
	return matches, err
}

func (r *repository) GetTeamsByIDs(ids []uint) (map[uint]string, error) {
	var teams []teammodels.Team
	if err := r.db.Where("id IN ?", ids).Find(&teams).Error; err != nil {
		return nil, err
	}
	out := map[uint]string{}
	for _, t := range teams {
		out[t.ID] = t.Name
	}
	return out, nil
}

func (r *repository) ListGoalsByMatch(matchID uint) ([]GoalRow, error) {
	var rows []GoalRow
	err := r.db.Table("goal_events g").
		Select("g.match_id as match_id, p.name as player_name").
		Joins("JOIN players p ON p.id = g.player_id").
		Where("g.match_id = ?", matchID).
		Scan(&rows).Error
	return rows, err
}
