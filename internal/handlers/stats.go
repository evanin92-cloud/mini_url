package handlers

import (
	"encoding/json"
	"net/http"

	"mini_url/internal/service"

	"github.com/go-chi/chi/v5"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "short_id")
	if shortID == "" {
		http.Error(w, "Short ID is required", http.StatusBadRequest)
		return
	}

	stats, err := h.statsService.GetStats(shortID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}