package constants

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleStaff  Role = "staff"
	RoleViewer Role = "viewer"
)

type PlayerPosition string

const (
	PositionForward    PlayerPosition = "penyerang"
	PositionMidfielder PlayerPosition = "gelandang"
	PositionDefender   PlayerPosition = "bertahan"
	PositionGoalkeeper PlayerPosition = "penjaga_gawang"
)

type MatchStatus string

const (
	MatchScheduled MatchStatus = "scheduled"
	MatchFinished  MatchStatus = "finished"
)
