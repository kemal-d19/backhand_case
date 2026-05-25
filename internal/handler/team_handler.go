package handler

import (
	"encoding/json"
	"insider_case/internal/model"
	"insider_case/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team model.Team

	err := json.NewDecoder(r.Body).Decode(&team)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTeam, err := h.teamService.CreateTeam(r.Context(), team)

	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTeam)
}

func (h *TeamHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.teamService.GetTeams(r.Context())
	if err != nil {
		http.Error(w, "Failed to get teams.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (h *TeamHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	// Expected URL: /teams/{id}
	// Example: /teams/3

	idStr := strings.TrimPrefix(r.URL.Path, "/team/get_by_id/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid team id", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.GetTeamByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := h.teamService.DeleteAllTeams(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
