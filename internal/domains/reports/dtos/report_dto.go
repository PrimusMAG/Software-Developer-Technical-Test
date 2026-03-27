package dtos

import "time"

type MatchReportItem struct {
	MatchID                 uint      `json:"matchId"`
	KickoffAt               time.Time `json:"kickoffAt"`
	HomeTeam                string    `json:"homeTeam"`
	AwayTeam                string    `json:"awayTeam"`
	FinalScore              string    `json:"finalScore"`
	FinalStatus             string    `json:"finalStatus"`
	TopScorer               string    `json:"topScorer"`
	HomeTeamWinAccumulation int       `json:"homeTeamWinAccumulation"`
	AwayTeamWinAccumulation int       `json:"awayTeamWinAccumulation"`
}
