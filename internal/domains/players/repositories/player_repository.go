package repositories

import (
	"gorm.io/gorm"

	playermodels "football-api/internal/domains/players/models"
)

type Repository interface {
	Create(player *playermodels.Player) error
	Update(player *playermodels.Player) error
	Delete(id uint) error
	GetByID(id uint) (*playermodels.Player, error)
	List(offset, limit int, search, position string, teamID uint) ([]playermodels.Player, int64, error)
}

type repository struct{ db *gorm.DB }

func New(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(player *playermodels.Player) error { return r.db.Create(player).Error }
func (r *repository) Update(player *playermodels.Player) error { return r.db.Save(player).Error }
func (r *repository) Delete(id uint) error                     { return r.db.Delete(&playermodels.Player{}, id).Error }

func (r *repository) GetByID(id uint) (*playermodels.Player, error) {
	var player playermodels.Player
	if err := r.db.First(&player, id).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *repository) List(offset, limit int, search, position string, teamID uint) ([]playermodels.Player, int64, error) {
	var players []playermodels.Player
	var total int64
	q := r.db.Model(&playermodels.Player{})
	if search != "" {
		q = q.Where("name ILIKE ?", "%"+search+"%")
	}
	if position != "" {
		q = q.Where("position = ?", position)
	}
	if teamID > 0 {
		q = q.Where("team_id = ?", teamID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset(offset).Limit(limit).Find(&players).Error; err != nil {
		return nil, 0, err
	}
	return players, total, nil
}
