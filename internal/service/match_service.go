package service

import (
	"context"
	"errors"
	"fmt"
	"insider_case/internal/model"
	"insider_case/internal/repository"
	"math/rand"
	"time"
)

type MatchService struct {
	matchRepository *repository.MatchRepository
	teamRepository  *repository.TeamRepository
}

func NewMatchService(matchRepository *repository.MatchRepository, teamRepository *repository.TeamRepository) *MatchService {
	return &MatchService{matchRepository: matchRepository, teamRepository: teamRepository}
}

func (s *MatchService) CreateMatch(ctx context.Context, match *model.Match) (*model.Match, error) {
	if match.WeekNumber <= 0 {
		return nil, errors.New("week number must be greater than 0")
	}

	if match.HomeTeamID <= 0 {
		return nil, errors.New("home team id must be greater than 0")
	}

	if match.AwayTeamID <= 0 {
		return nil, errors.New("away team id must be greater than 0")
	}

	if match.HomeTeamID == match.AwayTeamID {
		return nil, errors.New("home team and away team cannot be same")
	}

	if match.HomeScore < 0 || match.AwayScore < 0 {
		return nil, errors.New("scores cannot be negative")
	}

	return s.matchRepository.CreateMatch(ctx, match)
}

func (s *MatchService) GetMatchByID(ctx context.Context, id int) (*model.Match, error) {
	if id <= 0 {
		return nil, errors.New("match id must be greater than 0")
	}

	return s.matchRepository.GetMatchByID(ctx, id)
}

func (s *MatchService) GetAllMatches(ctx context.Context) ([]model.Match, error) {
	return s.matchRepository.GetAllMatches(ctx, false)
}

func (s *MatchService) GenerateMatches(ctx context.Context, total_week int) ([]model.Match, error) {
	teams, err := s.teamRepository.GetAllTeams(ctx)
	if err != nil {
		return nil, err
	}
	if len(teams) < 2 {
		return nil, fmt.Errorf("At least 2 teams are required to generate matches.")
	}

	matches := GenerateDoubleRoundMatches(teams, total_week)

	for i := 0; i < len(matches); i++ {
		_, err = s.matchRepository.CreateMatch(ctx, &matches[i])
	}

	matches, err = s.matchRepository.GetAllMatches(ctx, false)

	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchService) UpdateMatchResult(
	ctx context.Context,
	matchID int,
	homeScore int,
	awayScore int,
	isPlayed bool,
) (*model.Match, error) {
	if matchID <= 0 {
		return nil, fmt.Errorf("invalid match id")
	}

	if homeScore < 0 {
		return nil, fmt.Errorf("home score cannot be negative")
	}

	if awayScore < 0 {
		return nil, fmt.Errorf("away score cannot be negative")
	}

	updatedMatch, err := s.matchRepository.UpdateMatchResult(
		ctx,
		matchID,
		homeScore,
		awayScore,
		isPlayed,
	)
	if err != nil {
		return nil, err
	}

	return updatedMatch, nil
}

func GenerateDoubleRoundMatches(teams []model.Team, weekN int) []model.Match {
	matches := make([]model.Match, 0)

	n := len(teams)
	if n < 2 || weekN <= 0 {
		return matches
	}

	totalMatches := n * (n - 1) // double round-robin: each pair plays home and away

	weekList := make([]int, 0, totalMatches)

	baseRep := totalMatches / weekN
	extraRep := totalMatches % weekN

	for week := 1; week <= weekN; week++ {
		repeat := baseRep
		if week <= extraRep {
			repeat++
		}

		for i := 0; i < repeat; i++ {
			weekList = append(weekList, week)
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ {
			if len(weekList) == 0 {
				return matches
			}

			firstWeekIndex := rng.Intn(len(weekList))
			firstWeek := weekList[firstWeekIndex]
			weekList = removeByIndex(weekList, firstWeekIndex)

			secondWeek := firstWeek

			if len(weekList) > 0 {
				secondWeekIndex := rng.Intn(len(weekList))

				// Try to avoid putting both matches in the same week
				for attempts := 0; attempts < 20; attempts++ {
					if weekList[secondWeekIndex] != firstWeek {
						break
					}
					secondWeekIndex = rng.Intn(len(weekList))
				}

				secondWeek = weekList[secondWeekIndex]
				weekList = removeByIndex(weekList, secondWeekIndex)
			}

			firstMatch := model.Match{
				WeekNumber: firstWeek,
				HomeTeamID: teams[i].ID,
				AwayTeamID: teams[j].ID,
				HomeScore:  0,
				AwayScore:  0,
				IsPlayed:   false,
			}

			secondMatch := model.Match{
				WeekNumber: secondWeek,
				HomeTeamID: teams[j].ID,
				AwayTeamID: teams[i].ID,
				HomeScore:  0,
				AwayScore:  0,
				IsPlayed:   false,
			}

			matches = append(matches, firstMatch, secondMatch)
		}
	}

	return matches
}

func (s *MatchService) DeleteAll(ctx context.Context) error {
	err := s.matchRepository.DeleteAllMatches(ctx)
	if err != nil {
		return err
	}
	return nil
}

func removeByIndex(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}
