package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"insider_case/internal/model"
	"insider_case/internal/service"
)

type UpdateMatchResultRequest struct {
	HomeScore int  `json:"home_score"`
	AwayScore int  `json:"away_score"`
	IsPlayed  bool `json:"is_played"`
}

type MatchHandler struct {
	matchService *service.MatchService
}

func NewMatchHandler(matchService *service.MatchService) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

func (h *MatchHandler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var match model.Match

	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	createdMatch, err := h.matchService.CreateMatch(r.Context(), &match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(createdMatch)
}

func (h *MatchHandler) GetMatchByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/match/get_by_id/")
	idStr = strings.TrimSpace(idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid match id", http.StatusBadRequest)
		return
	}

	match, err := h.matchService.GetMatchByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(match)
}

func (h *MatchHandler) UpdateMatchResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/match/update_result/")
	matchID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid match id", http.StatusBadRequest)
		return
	}

	var req UpdateMatchResultRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedMatch, err := h.matchService.UpdateMatchResult(
		r.Context(),
		matchID,
		req.HomeScore,
		req.AwayScore,
		req.IsPlayed,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(updatedMatch)
}

func (h *MatchHandler) GetAllMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	matches, err := h.matchService.GetAllMatches(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(matches)
}

func (h *MatchHandler) GenerateMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TotalWeek int `json:"total_week"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.TotalWeek <= 0 {
		http.Error(w, "total_week must be greater than 0", http.StatusBadRequest)
		return
	}

	matches, err := h.matchService.GenerateMatches(r.Context(), req.TotalWeek)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(matches)
}

func (h *MatchHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := h.matchService.DeleteAll(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
