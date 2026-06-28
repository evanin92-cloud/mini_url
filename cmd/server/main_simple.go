package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"mini_url/internal/handlers"
	"mini_url/internal/repository"
	"mini_url/internal/service"
)

func main() {
	// Используем хранилище в памяти (без БД)
	linkRepo := repository.NewMemoryLinkRepository()

	linkService := service.NewLinkService(linkRepo)
	linkHandler := handlers.NewLinkHandler(linkService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", linkHandler.Index)
	r.Get("/{short_id}", linkHandler.Redirect)
	r.Post("/api/v1/shorten", linkHandler.Create)

	log.Printf("🚀 Server starting on port 8080")
	log.Printf("📍 Visit http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}