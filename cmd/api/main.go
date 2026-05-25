package main

import (
	"context"
	"fmt"
	"insider_case/internal/database"
	"insider_case/internal/handler"
	"insider_case/internal/repository"
	"insider_case/internal/service"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	ctx := context.Background()

	db, err := database.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	dbcon := database.NewDataBase(db)
	dbcon.CreateTables(ctx)

	teamRepository := repository.NewTeamRepository(db)
	teamService := service.NewTeamService(teamRepository)
	teamHandler := handler.NewTeamHandler(teamService)

	matchRepository := repository.NewMatchRepository(db)
	matchService := service.NewMatchService(matchRepository, teamRepository)
	matchHandler := handler.NewMatchHandler(matchService)

	leagueService := service.NewLeqgueService(matchRepository, teamRepository)
	legueHandler := handler.NewLeagueHandler(leagueService)

	mux.HandleFunc("POST /team/create/", teamHandler.CreateTeam)
	mux.HandleFunc("/team/get_all", teamHandler.GetTeams)
	mux.HandleFunc("/team/delete_all/", teamHandler.DeleteAll)

	mux.HandleFunc("/match/generate/", matchHandler.GenerateMatches)
	mux.HandleFunc("/match/get_all", matchHandler.GetAllMatches)
	mux.HandleFunc("/match/delete_all", matchHandler.DeleteAll)

	mux.HandleFunc("/league/play_current_week/", legueHandler.PlayCurrentWeek)
	mux.HandleFunc("/league/play_all_weeks/", legueHandler.PlayAlltWeek)
	mux.HandleFunc("/league/get_table/", legueHandler.GetResultTable)
	mux.HandleFunc("/league/simulate/", legueHandler.GetPredictionTable)

	fmt.Println("Server running on port 8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
