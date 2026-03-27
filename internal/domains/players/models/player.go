package models

import (
	"time"

	team "football-api/internal/domains/teams/models"

	"gorm.io/gorm"
)

type Player struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TeamID       uint           `gorm:"not null;index;uniqueIndex:idx_team_jersey" json:"teamId"`
	Team         team.Team      `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Name         string         `gorm:"size:120;not null" json:"name"`
	HeightCM     int            `gorm:"not null" json:"heightCm"`
	WeightKG     int            `gorm:"not null" json:"weightKg"`
	Position     string         `gorm:"size:30;index;not null" json:"position"`
	JerseyNumber int            `gorm:"not null;uniqueIndex:idx_team_jersey" json:"jerseyNumber"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
