package dtos

type CreateMatchRequest struct {
	KickoffAt  string `json:"kickoffAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	HomeTeamID uint   `json:"homeTeamId" validate:"required,gt=0"`
	AwayTeamID uint   `json:"awayTeamId" validate:"required,gt=0,nefield=HomeTeamID"`
}

type GoalInput struct {
	TeamID   uint `json:"teamId" validate:"required,gt=0"`
	PlayerID uint `json:"playerId" validate:"required,gt=0"`
	Minute   int  `json:"minute" validate:"required,gte=0,lte=130"`
}

type SubmitResultRequest struct {
	HomeScore int         `json:"homeScore" validate:"gte=0,lte=50"`
	AwayScore int         `json:"awayScore" validate:"gte=0,lte=50"`
	Goals     []GoalInput `json:"goals" validate:"dive"`
}
