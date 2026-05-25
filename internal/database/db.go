package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func NewPostgresConnection(ctx context.Context) (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbSSLMode,
	)

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

type DataBase struct {
	db *pgxpool.Pool
}

func NewDataBase(db *pgxpool.Pool) *DataBase {
	return &DataBase{db: db}
}

func (r *DataBase) CreateTables(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS teams (
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
                );`
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("Failed to create tables: %w", err)
	}
	log.Println("Database tables initialized Successfully")
	return nil
}
