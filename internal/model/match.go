package model

type Match struct {
	ID         int  `json:"id" db:"id"`
	WeekNumber int  `json:"week_number" db:"week_number"`
	HomeTeamID int  `json:"home_team_id" db:"home_team_id"`
	AwayTeamID int  `json:"away_team_id" db:"away_team_id"`
	HomeScore  int  `json:"home_score" db:"home_score"`
	AwayScore  int  `json:"away_score" db:"away_score"`
	IsPlayed   bool `json:"is_played" db:"is_played"`
}
