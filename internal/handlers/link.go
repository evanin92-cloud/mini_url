package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"mini_url/internal/middleware"
	"mini_url/internal/service"

	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	linkService *service.LinkService
}

func NewLinkHandler(linkService *service.LinkService) *LinkHandler {
	return &LinkHandler{linkService: linkService}
}

type ShortenRequest struct {
	URL      string `json:"url"`
	CustomID string `json:"custom_id,omitempty"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	Original string `json:"original"`
	Error    string `json:"error,omitempty"`
}

func (h *LinkHandler) Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (h *LinkHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	tmpl, err := template.ParseFiles("web/templates/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	links, _ := h.linkService.GetLinksByUser(userID)
	data := map[string]interface{}{
		"Links": links,
	}
	tmpl.Execute(w, data)
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	link, err := h.linkService.CreateLink(req.URL, userID, req.CustomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ShortenResponse{
		ShortURL: "http://localhost:8080/" + link.ShortID,
		Original: link.OriginalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *LinkHandler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req struct {
		URLs []string `json:"urls"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	links, err := h.linkService.CreateBatch(req.URLs, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func (h *LinkHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	links, err := h.linkService.GetLinksByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func (h *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	linkID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid link ID", http.StatusBadRequest)
		return
	}

	if err := h.linkService.DeleteLink(linkID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "short_id")
	if shortID == "" {
		http.NotFound(w, r)
		return
	}

	originalURL, err := h.linkService.Redirect(shortID)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *LinkHandler) AdminPanel(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/admin.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}