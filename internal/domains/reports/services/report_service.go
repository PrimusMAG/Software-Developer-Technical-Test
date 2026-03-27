package services

import (
	"fmt"
	"sync"
	"time"

	"football-api/internal/domains/reports/dtos"
	"football-api/internal/domains/reports/repositories"
)

type Service interface {
	List() ([]dtos.MatchReportItem, error)
	Revalidate()
}

type cacheState struct {
	data      []dtos.MatchReportItem
	expiresAt time.Time
}

type service struct {
	repo repositories.Repository
	mu   sync.RWMutex
	c    cacheState
}

func New(repo repositories.Repository) Service {
	return &service{repo: repo}
}

func (s *service) List() ([]dtos.MatchReportItem, error) {
	s.mu.RLock()
	if time.Now().Before(s.c.expiresAt) && len(s.c.data) > 0 {
		defer s.mu.RUnlock()
		return s.c.data, nil
	}
	s.mu.RUnlock()

	matches, err := s.repo.ListFinishedMatches()
	if err != nil {
		return nil, err
	}

	teamIDsSet := map[uint]bool{}
	for _, m := range matches {
		teamIDsSet[m.HomeTeamID] = true
		teamIDsSet[m.AwayTeamID] = true
	}
	teamIDs := make([]uint, 0, len(teamIDsSet))
	for id := range teamIDsSet {
		teamIDs = append(teamIDs, id)
	}
	teams, err := s.repo.GetTeamsByIDs(teamIDs)
	if err != nil {
		return nil, err
	}

	winCount := map[uint]int{}
	scorerCount := map[string]int{}
	report := make([]dtos.MatchReportItem, 0, len(matches))
	for _, m := range matches {
		goals, err := s.repo.ListGoalsByMatch(m.ID)
		if err != nil {
			return nil, err
		}
		for _, g := range goals {
			scorerCount[g.PlayerName]++
		}
		topScorer := "-"
		maxGoal := 0
		for name, total := range scorerCount {
			if total > maxGoal {
				maxGoal = total
				topScorer = name
			}
		}
		finalStatus := "Draw"
		if m.HomeScore > m.AwayScore {
			finalStatus = "Tim Home Menang"
			winCount[m.HomeTeamID]++
		} else if m.AwayScore > m.HomeScore {
			finalStatus = "Tim Away Menang"
			winCount[m.AwayTeamID]++
		}
		report = append(report, dtos.MatchReportItem{
			MatchID:                 m.ID,
			KickoffAt:               m.KickoffAt,
			HomeTeam:                teams[m.HomeTeamID],
			AwayTeam:                teams[m.AwayTeamID],
			FinalScore:              fmt.Sprintf("%d-%d", m.HomeScore, m.AwayScore),
			FinalStatus:             finalStatus,
			TopScorer:               topScorer,
			HomeTeamWinAccumulation: winCount[m.HomeTeamID],
			AwayTeamWinAccumulation: winCount[m.AwayTeamID],
		})
	}

	s.mu.Lock()
	s.c = cacheState{data: report, expiresAt: time.Now().Add(60 * time.Second)}
	s.mu.Unlock()
	return report, nil
}

func (s *service) Revalidate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.c = cacheState{}
}
