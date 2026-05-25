package service

import (
	"context"
	"errors"
	"fmt"
	"insider_case/internal/model"
	"insider_case/internal/repository"
)

// type TeamService interface {
// 	CreateTeam(ctx context.Context, team model.Team) (*model.Team, error)
// 	GetTeams(ctx context.Context) ([]model.Team, error)
// }

type TeamService struct {
	teamRepository *repository.TeamRepository
}

func NewTeamService(teamRepository *repository.TeamRepository) *TeamService {
	return &TeamService{teamRepository: teamRepository}
}

func (s *TeamService) CreateTeam(ctx context.Context, team model.Team) (*model.Team, error) {
	if team.Name == "" {
		return nil, errors.New("Team name is required.")
	}
	if team.Strength <= 0 {
		return nil, errors.New("Team strength is required.")
	}
	return s.teamRepository.Create(ctx, team)
}

func (s *TeamService) GetTeams(ctx context.Context) ([]model.Team, error) {
	return s.teamRepository.GetAllTeams(ctx)
}

func (s *TeamService) GetTeamByID(ctx context.Context, id int) (*model.Team, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid team id")
	}

	team, err := s.teamRepository.GetTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) DeleteAllTeams(ctx context.Context) error {
	err := s.teamRepository.DeleteAll(ctx)

	if err != nil {
		return err
	}
	return nil
}
