package repositories

import (
	"gorm.io/gorm"

	usermodels "football-api/internal/domains/users/models"
)

type Repository interface {
	FindByEmail(email string) (*usermodels.User, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(email string) (*usermodels.User, error) {
	var user usermodels.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
