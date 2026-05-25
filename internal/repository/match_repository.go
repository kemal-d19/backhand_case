package repository

import (
	"context"
	"fmt"

	"insider_case/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository struct {
	db *pgxpool.Pool
}

func NewMatchRepository(db *pgxpool.Pool) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) CreateMatch(ctx context.Context, match *model.Match) (*model.Match, error) {
	query := `
		INSERT INTO matches 
		(week_number, home_team_id, away_team_id, home_score, away_score, is_played)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, week_number, home_team_id, away_team_id, home_score, away_score, is_played;
	`

	var createdMatch model.Match

	err := r.db.QueryRow(
		ctx,
		query,
		match.WeekNumber,
		match.HomeTeamID,
		match.AwayTeamID,
		match.HomeScore,
		match.AwayScore,
		match.IsPlayed,
	).Scan(
		&createdMatch.ID,
		&createdMatch.WeekNumber,
		&createdMatch.HomeTeamID,
		&createdMatch.AwayTeamID,
		&createdMatch.HomeScore,
		&createdMatch.AwayScore,
		&createdMatch.IsPlayed,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create match: %w", err)
	}

	return &createdMatch, nil
}

func (r *MatchRepository) GetMatchByID(ctx context.Context, id int) (*model.Match, error) {
	query := `
		SELECT id, week_number, home_team_id, away_team_id, home_score, away_score, is_played
		FROM matches
		WHERE id = $1;
	`

	var match model.Match

	err := r.db.QueryRow(ctx, query, id).Scan(
		&match.ID,
		&match.WeekNumber,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.HomeScore,
		&match.AwayScore,
		&match.IsPlayed,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get match by id: %w", err)
	}

	return &match, nil
}

func (r *MatchRepository) GetAllMatches(ctx context.Context, is_played bool) ([]model.Match, error) {
	query := `
		SELECT id, week_number, home_team_id, away_team_id, home_score, away_score, is_played
		FROM matches
		ORDER BY week_number, id;
	`

	if is_played == true {
		query = `
		SELECT id, week_number, home_team_id, away_team_id, home_score, away_score, is_played
		FROM matches WHERE is_played = true
		ORDER BY week_number, id;
	`

	}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}
	defer rows.Close()

	var matches []model.Match

	for rows.Next() {
		var match model.Match

		err := rows.Scan(
			&match.ID,
			&match.WeekNumber,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.HomeScore,
			&match.AwayScore,
			&match.IsPlayed,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return matches, nil
}

func (r MatchRepository) DeleteAllMatches(ctx context.Context) error {
	query := `DELETE FROM matches;`

	_, err := r.db.Exec(ctx, query)

	if err != nil {
		return err
	}
	return err
}

func (r *MatchRepository) UpdateMatchResult(
	ctx context.Context,
	matchID int,
	homeScore int,
	awayScore int,
	isPlayed bool,
) (*model.Match, error) {
	query := `
		UPDATE matches
		SET 
			home_score = $1,
			away_score = $2,
			is_played = $3
		WHERE id = $4
		RETURNING id, week_number, home_team_id, away_team_id, home_score, away_score, is_played;
	`

	var match model.Match

	err := r.db.QueryRow(
		ctx,
		query,
		homeScore,
		awayScore,
		isPlayed,
		matchID,
	).Scan(
		&match.ID,
		&match.WeekNumber,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.HomeScore,
		&match.AwayScore,
		&match.IsPlayed,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update match result for match id %d: %w", matchID, err)
	}

	return &match, nil
}
