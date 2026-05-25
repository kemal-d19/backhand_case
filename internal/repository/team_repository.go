package repository

import (
	"context"
	"fmt"
	"insider_case/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team model.Team) (*model.Team, error) {
	query := `
	INSERT INTO teams (name, strength)
	VALUES ($1, $2)
	RETURNING id, name, strength`
	var createTeam model.Team

	err := r.db.QueryRow(ctx, query, team.Name, team.Strength).Scan(
		&createTeam.ID,
		&createTeam.Name,
		&createTeam.Strength,
	)
	if err != nil {
		return nil, err
	}
	return &createTeam, nil

}

func (r *TeamRepository) GetAllTeams(ctx context.Context) ([]model.Team, error) {
	query := `SELECT id, name, strength FROM teams`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := make([]model.Team, 0)
	for rows.Next() {
		var team model.Team
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Strength,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

// GetTeamByID fetches a single team by its unique ID.
func (r *TeamRepository) GetTeamByID(ctx context.Context, id int) (*model.Team, error) {
	query := `SELECT id, name, strength FROM teams WHERE id = $1;`

	var team model.Team
	err := r.db.QueryRow(ctx, query, id).Scan(&team.ID, &team.Name, &team.Strength)
	if err != nil {
		return nil, fmt.Errorf("failed to get team by id %d: %w", id, err)
	}

	return &team, nil
}

func (r *TeamRepository) DeleteAll(ctx context.Context) error {
	query := `DELETE FROM teams;`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return err
	}
	return err
}
