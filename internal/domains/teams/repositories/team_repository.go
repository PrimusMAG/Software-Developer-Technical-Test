package repositories

import (
	"gorm.io/gorm"

	teammodels "football-api/internal/domains/teams/models"
)

type Repository interface {
	Create(team *teammodels.Team) error
	Update(team *teammodels.Team) error
	Delete(id uint) error
	GetByID(id uint) (*teammodels.Team, error)
	List(offset, limit int, search, city string) ([]teammodels.Team, int64, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(team *teammodels.Team) error { return r.db.Create(team).Error }
func (r *repository) Update(team *teammodels.Team) error { return r.db.Save(team).Error }
func (r *repository) Delete(id uint) error               { return r.db.Delete(&teammodels.Team{}, id).Error }

func (r *repository) GetByID(id uint) (*teammodels.Team, error) {
	var team teammodels.Team
	if err := r.db.First(&team, id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *repository) List(offset, limit int, search, city string) ([]teammodels.Team, int64, error) {
	var teams []teammodels.Team
	var total int64
	q := r.db.Model(&teammodels.Team{})
	if search != "" {
		q = q.Where("name ILIKE ?", "%"+search+"%")
	}
	if city != "" {
		q = q.Where("hq_city ILIKE ?", "%"+city+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		return nil, 0, err
	}
	return teams, total, nil
}
