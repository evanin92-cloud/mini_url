package handlers

import (
	"net/http"

	"mini_url/internal/service"

	"github.com/go-chi/chi/v5"
)

// RedirectHandler обрабатывает редирект по короткой ссылке
type RedirectHandler struct {
	linkService *service.LinkService
}

func NewRedirectHandler(linkService *service.LinkService) *RedirectHandler {
	return &RedirectHandler{linkService: linkService}
}

// Redirect перенаправляет пользователя по короткой ссылке на оригинальный URL
func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
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