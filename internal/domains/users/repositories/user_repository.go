package repositories

import (
	"gorm.io/gorm"

	usermodels "football-api/internal/domains/users/models"
)

type Repository interface {
	Create(user *usermodels.User) error
	List(offset, limit int, search string) ([]usermodels.User, int64, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *usermodels.User) error {
	return r.db.Create(user).Error
}

func (r *repository) List(offset, limit int, search string) ([]usermodels.User, int64, error) {
	var users []usermodels.User
	var total int64
	q := r.db.Model(&usermodels.User{})
	if search != "" {
		q = q.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
