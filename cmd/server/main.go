package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"mini_url/internal/config"
	"mini_url/internal/handlers"
	authMiddleware "mini_url/internal/middleware"
	"mini_url/internal/repository"
	"mini_url/internal/service"
	"mini_url/pkg/database"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к PostgreSQL
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Подключение к Redis
	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	linkRepo := repository.NewLinkRepository(db, redisClient)
	statsRepo := repository.NewStatsRepository(db)

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	linkService := service.NewLinkService(linkRepo)
	statsService := service.NewStatsService(statsRepo)

	// Инициализация хендлеров
	authHandler := handlers.NewAuthHandler(authService)
	linkHandler := handlers.NewLinkHandler(linkService)
	statsHandler := handlers.NewStatsHandler(statsService)

	// Настройка роутера
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "gzip"))

	// Статические файлы
	workDir, _ := os.Getwd()
	filesDir := http.Dir(workDir + "/web/static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(filesDir)))

	// Публичные маршруты
	r.Get("/", linkHandler.Index)
	r.Get("/{short_id}", linkHandler.Redirect)
	r.Post("/api/v1/register", authHandler.Register)
	r.Post("/api/v1/login", authHandler.Login)

	// Защищённые маршруты (требуют JWT)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware(cfg.JWTSecret))

		r.Get("/dashboard", linkHandler.Dashboard)
		r.Post("/api/v1/shorten", linkHandler.Create)
		r.Post("/api/v1/shorten/batch", linkHandler.CreateBatch)
		r.Get("/api/v1/links", linkHandler.GetLinks)
		r.Delete("/api/v1/links/{id}", linkHandler.Delete)
		r.Get("/api/v1/stats/{short_id}", statsHandler.GetStats)

		// Административные маршруты
		r.Get("/admin", linkHandler.AdminPanel)
		r.Get("/api/v1/admin/users", authHandler.GetUsers)
		r.Delete("/api/v1/admin/users/{id}", authHandler.DeleteUser)
	})

	// Запуск сервера
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Printf("📍 Visit http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}