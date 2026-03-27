package repositories

import (
	"time"

	"gorm.io/gorm"

	matchmodels "football-api/internal/domains/matches/models"
	playermodels "football-api/internal/domains/players/models"
)

type Repository interface {
	Create(match *matchmodels.Match) error
	Update(match *matchmodels.Match) error
	Delete(id uint) error
	GetByID(id uint) (*matchmodels.Match, error)
	List(offset, limit int, status string, homeTeamID, awayTeamID uint, dateFrom, dateTo *time.Time) ([]matchmodels.Match, int64, error)
	GetPlayer(id uint) (*playermodels.Player, error)
	CreateGoals(tx *gorm.DB, goals []matchmodels.GoalEvent) error
	DeleteGoalsByMatch(tx *gorm.DB, matchID uint) error
	WithTx(fn func(tx *gorm.DB) error) error
}

type repository struct{ db *gorm.DB }

func New(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(match *matchmodels.Match) error { return r.db.Create(match).Error }
func (r *repository) Update(match *matchmodels.Match) error { return r.db.Save(match).Error }
func (r *repository) Delete(id uint) error                  { return r.db.Delete(&matchmodels.Match{}, id).Error }

func (r *repository) GetByID(id uint) (*matchmodels.Match, error) {
	var match matchmodels.Match
	if err := r.db.First(&match, id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *repository) List(offset, limit int, status string, homeTeamID, awayTeamID uint, dateFrom, dateTo *time.Time) ([]matchmodels.Match, int64, error) {
	var matches []matchmodels.Match
	var total int64
	q := r.db.Model(&matchmodels.Match{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if homeTeamID > 0 {
		q = q.Where("home_team_id = ?", homeTeamID)
	}
	if awayTeamID > 0 {
		q = q.Where("away_team_id = ?", awayTeamID)
	}
	if dateFrom != nil {
		q = q.Where("kickoff_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		q = q.Where("kickoff_at <= ?", *dateTo)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("kickoff_at ASC").Offset(offset).Limit(limit).Find(&matches).Error; err != nil {
		return nil, 0, err
	}
	return matches, total, nil
}

func (r *repository) GetPlayer(id uint) (*playermodels.Player, error) {
	var p playermodels.Player
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) CreateGoals(tx *gorm.DB, goals []matchmodels.GoalEvent) error {
	if len(goals) == 0 {
		return nil
	}
	return tx.Create(&goals).Error
}

func (r *repository) DeleteGoalsByMatch(tx *gorm.DB, matchID uint) error {
	return tx.Where("match_id = ?", matchID).Delete(&matchmodels.GoalEvent{}).Error
}

func (r *repository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
