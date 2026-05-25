package service

import (
	"context"
	"fmt"
	"insider_case/internal/model"
	"insider_case/internal/repository"
	"math"
	"math/rand"
	"sort"
)

type LeagueService struct {
	matchRepo *repository.MatchRepository
	teamRepo  *repository.TeamRepository
}

func NewLeqgueService(matchRepo *repository.MatchRepository, teamRepo *repository.TeamRepository) *LeagueService {
	return &LeagueService{matchRepo: matchRepo, teamRepo: teamRepo}
}

func (s *LeagueService) GetResultTable(ctx context.Context) ([]model.LeagueTableRow, error) {

	teams, _ := s.teamRepo.GetAllTeams(ctx)
	played_matches, _ := s.matchRepo.GetAllMatches(ctx, true)

	table := CalculateTable(played_matches, teams)
	return table, nil
}

func (s *LeagueService) PlayCurrentWeek(ctx context.Context, unpred_coef float32) ([]model.Match, error) {
	allMatches, err := s.matchRepo.GetAllMatches(ctx, false)
	teams, err := s.teamRepo.GetAllTeams(ctx)
	if err != nil {
		return nil, err
	}
	current_week := 1
	max_week := 1
	for i := 0; i < len(allMatches); i++ {
		if allMatches[i].WeekNumber > max_week {
			max_week = allMatches[i].WeekNumber
		}
		if allMatches[i].IsPlayed == true && allMatches[i].WeekNumber+1 > current_week {
			current_week = allMatches[i].WeekNumber + 1
		}
	}

	if current_week == max_week+1 {
		return nil, fmt.Errorf("All matches are played")
	}

	curretMatches := make([]model.Match, 0)
	for i := 0; i < len(allMatches); i++ {
		if allMatches[i].WeekNumber == current_week {
			curretMatches = append(curretMatches, allMatches[i])
		}
	}

	newMatches := make([]model.Match, 0)
	for i := 0; i < len(curretMatches); i++ {

		home_goal, away_goal := simulateMatches(
			StrgByID(teams, curretMatches[i].HomeTeamID),
			StrgByID(teams, curretMatches[i].AwayTeamID),
			unpred_coef,
		)
		newMatc, _ := s.matchRepo.UpdateMatchResult(ctx, curretMatches[i].ID, home_goal, away_goal, true)
		newMatches = append(newMatches, *newMatc)
	}

	return newMatches, nil
}

func (s *LeagueService) PlayAllWeeks(ctx context.Context, unpred_coef float32) ([]model.Match, error) {
	allMatches, err := s.matchRepo.GetAllMatches(ctx, false)
	teams, err := s.teamRepo.GetAllTeams(ctx)
	if err != nil {
		return nil, err
	}

	max_week := 1
	current_week := 1
	unplayedMatches := make([]model.Match, 0)
	for i := 0; i < len(allMatches); i++ {
		if allMatches[i].IsPlayed == false {
			unplayedMatches = append(unplayedMatches, allMatches[i])
		}
		if allMatches[i].WeekNumber > max_week {
			max_week = allMatches[i].WeekNumber
		}
		if allMatches[i].IsPlayed == true && allMatches[i].WeekNumber+1 > current_week {
			current_week = allMatches[i].WeekNumber + 1
		}
	}

	if current_week == max_week+1 {
		return nil, fmt.Errorf("All matches are played")
	}

	newMatches := make([]model.Match, 0)
	for i := 0; i < len(unplayedMatches); i++ {
		home_goal, away_goal := simulateMatches(
			StrgByID(teams, unplayedMatches[i].HomeTeamID),
			StrgByID(teams, unplayedMatches[i].AwayTeamID),
			unpred_coef,
		)

		newMatc, _ := s.matchRepo.UpdateMatchResult(ctx, unplayedMatches[i].ID, home_goal, away_goal, true)
		newMatches = append(newMatches, *newMatc)

	}

	return newMatches, nil
}

func (s *LeagueService) MonteCarloSimulation(ctx context.Context, unpred_coef float32) ([]model.PredictionTableRow, error) {
	allMatches, err := s.matchRepo.GetAllMatches(ctx, false)
	teams, err := s.teamRepo.GetAllTeams(ctx)
	if err != nil {
		return nil, err
	}
	current_week := 1
	max_week := 1
	playedMatches := make([]model.Match, 0)

	for i := 0; i < len(allMatches); i++ {
		if allMatches[i].IsPlayed == true {
			playedMatches = append(playedMatches, allMatches[i])
		}
		if allMatches[i].WeekNumber > max_week {
			max_week = allMatches[i].WeekNumber
		}
		if allMatches[i].IsPlayed == true && allMatches[i].WeekNumber+1 > current_week {
			current_week = allMatches[i].WeekNumber + 1
		}
	}

	if current_week == max_week+1 {
		return nil, fmt.Errorf("All matches are played")
	}
	type WinCount struct {
		Team  string
		Count int
	}
	WinCounts := make([]WinCount, 0)
	for _, v := range teams {
		w := WinCount{Team: v.Name, Count: 0}
		WinCounts = append(WinCounts, w)
	}

	N_ITER := 100000

	for n := 0; n < N_ITER; n++ {

		// Start each simulation with already played matches
		allPlayedMatches := make([]model.Match, 0)
		allPlayedMatches = append(allPlayedMatches, playedMatches...)

		for w := current_week; w <= max_week; w++ {

			currentMatches := make([]model.Match, 0)

			for i := 0; i < len(allMatches); i++ {
				if allMatches[i].WeekNumber == w {
					currentMatches = append(currentMatches, allMatches[i])
				}
			}

			matchResults := PlayMatches(currentMatches, teams, unpred_coef)

			// Add simulated results to this simulation
			allPlayedMatches = append(allPlayedMatches, matchResults...)
		}
		// Calculate table using played + simulated matches
		resultTable := CalculateTable(allPlayedMatches, teams)

		for k := 0; k < 3 && k < len(resultTable); k++ {
			fmt.Printf("Team: %s, points: %d\n", resultTable[k].Team, resultTable[k].Points)
		}

		fmt.Println("----------------------------------")

		for j := 0; j < len(WinCounts); j++ {
			if WinCounts[j].Team == resultTable[0].Team {
				WinCounts[j].Count++
				break
			}
		}
	}

	total := 0
	for _, v := range WinCounts {
		total += v.Count
	}

	predTable := make([]model.PredictionTableRow, 0)

	for _, v := range WinCounts {
		pred_row := model.PredictionTableRow{
			Team:       v.Team,
			Prediction: float32(float32(100) * float32(v.Count) / float32(total)),
		}
		predTable = append(predTable, pred_row)
	}

	sort.Slice(predTable, func(i, j int) bool {
		return predTable[i].Prediction > predTable[j].Prediction
	})

	return predTable, nil
}

func CalculateTable(matches []model.Match, teams []model.Team) []model.LeagueTableRow {
	table := make([]model.LeagueTableRow, 0)

	for i := 0; i < len(teams); i++ {
		pts, played, won, drawn, lost, gd := 0, 0, 0, 0, 0, 0
		for j := 0; j < len(matches); j++ {
			if matches[j].HomeTeamID == teams[i].ID {
				s_home := matches[j].HomeScore
				s_away := matches[j].AwayScore

				played++
				gd = s_home - s_away
				if s_home > s_away {
					won++
					pts += 3
				} else if s_home < s_away {
					lost++
				} else {
					drawn++
					pts++
				}
			} else if matches[j].AwayTeamID == teams[i].ID {
				s_home := matches[j].HomeScore
				s_away := matches[j].AwayScore
				played++
				gd = s_away - s_home
				if s_away > s_home {
					won++
					pts += 3
				} else if s_away < s_home {
					lost++
				} else {
					drawn++
					pts++
				}

			}
		}
		row := model.LeagueTableRow{
			Team:   teams[i].Name,
			Points: pts,
			Played: played,
			Won:    won,
			Drawn:  drawn,
			Lost:   lost,
			GD:     gd,
		}
		table = append(table, row)
	}

	sort.Slice(table, func(i, j int) bool {
		if table[i].Points != table[j].Points {
			return table[i].Points > table[j].Points
		}

		return table[i].GD > table[j].GD
	})

	return table
}

func PlayMatches(matches []model.Match, teams []model.Team, unpred_coef float32) []model.Match {
	newMatches := make([]model.Match, 0)
	for i := 0; i < len(matches); i++ {
		home_goal, away_goal := simulateMatches(
			StrgByID(teams, matches[i].HomeTeamID),
			StrgByID(teams, matches[i].AwayTeamID),
			unpred_coef,
		)
		newMatch := model.Match{
			ID:         matches[i].ID,
			HomeTeamID: matches[i].HomeTeamID,
			AwayTeamID: matches[i].AwayTeamID,
			HomeScore:  home_goal,
			AwayScore:  away_goal,
			IsPlayed:   matches[i].IsPlayed,
		}

		newMatches = append(newMatches, newMatch)
	}

	return newMatches
}

func simulateMatches(home_strength float32, away_strength float32, uncertanity float32) (int, int) {
	var baseGoalRate float32 = 2.0
	var homeAdvantage float32 = 1.15

	homeNoise := (rand.Float32() * 2 * uncertanity) - uncertanity
	awaynoise := (rand.Float32() * 2 * uncertanity) - uncertanity

	lambda_home := (baseGoalRate * home_strength * homeAdvantage) + homeNoise
	if lambda_home < 0.5 {
		lambda_home = 0.5
	}

	lambda_away := (baseGoalRate * away_strength) + awaynoise
	if lambda_away < 0.5 {
		lambda_away = 0.5
	}

	homeGoals := poissonRandom(float64(lambda_home))
	awayGoals := poissonRandom(float64(lambda_away))

	return homeGoals, awayGoals
}

func poissonRandom(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= rand.Float64()
	}
	return k - 1
}

func StrgByID(teams []model.Team, id int) float32 {
	for _, team := range teams {
		if team.ID == id {
			return team.Strength
		}
	}
	return 0 // default value if not found
}
