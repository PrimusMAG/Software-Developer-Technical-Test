package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:120;uniqueIndex;not null" json:"name"`
	FoundedYear int            `gorm:"not null" json:"foundedYear"`
	HQAddress   string         `gorm:"size:255;not null" json:"hqAddress"`
	HQCity      string         `gorm:"size:120;index;not null" json:"hqCity"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
