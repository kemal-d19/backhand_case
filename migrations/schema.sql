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