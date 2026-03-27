package models

import (
	"time"

	player "football-api/internal/domains/players/models"
	team "football-api/internal/domains/teams/models"

	"gorm.io/gorm"
)

type Match struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	KickoffAt   time.Time      `gorm:"index;not null" json:"kickoffAt"`
	HomeTeamID  uint           `gorm:"not null;index" json:"homeTeamId"`
	HomeTeam    team.Team      `gorm:"foreignKey:HomeTeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	AwayTeamID  uint           `gorm:"not null;index" json:"awayTeamId"`
	AwayTeam    team.Team      `gorm:"foreignKey:AwayTeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Status      string         `gorm:"size:20;index;not null;default:scheduled" json:"status"`
	HomeScore   int            `gorm:"not null;default:0" json:"homeScore"`
	AwayScore   int            `gorm:"not null;default:0" json:"awayScore"`
	SubmittedAt *time.Time     `json:"submittedAt,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type GoalEvent struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	MatchID   uint          `gorm:"not null;index" json:"matchId"`
	Match     Match         `gorm:"foreignKey:MatchID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	TeamID    uint          `gorm:"not null;index" json:"teamId"`
	Team      team.Team     `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	PlayerID  uint          `gorm:"not null;index" json:"playerId"`
	Player    player.Player `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Minute    int           `gorm:"not null" json:"minute"`
	CreatedAt time.Time     `json:"createdAt"`
}
