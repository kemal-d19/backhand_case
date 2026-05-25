CREATE TABLE IF NOT EXISTS teams (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL UNIQUE,
				strength REAL NOT NULL
			    );

CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    week_number INT NOT NULL,
    home_team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    away_team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    home_score INT DEFAULT 0,
    away_score INT DEFAULT 0,
    is_played BOOLEAN DEFAULT FALSE
    );

-- Basic Queries


--Teams Table
INSERT INTO teams (name, strength) VALUES ($1, $2) RETURNING id, name, strength

SELECT id, name, strength FROM teams

SELECT id, name, strength FROM teams WHERE id = $1;

DELETE FROM teams;

--- Matches table
INSERT INTO matches 
	(week_number, home_team_id, away_team_id, home_score, away_score, is_played)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, week_number, home_team_id, away_team_id, home_score, away_score, is_played;

SELECT id, week_number, home_team_id, away_team_id, home_score, away_score, is_played FROM matches WHERE id = $1;

SELECT id, week_number, home_team_id, away_team_id, home_score, away_score, is_played FROM matches ORDER BY week_number, id;

UPDATE matches
    SET 
        home_score = $1,
        away_score = $2,
        is_played = $3
    WHERE id = $4
    RETURNING id, week_number, home_team_id, away_team_id, home_score, away_score, is_played;

DELETE FROM matches;