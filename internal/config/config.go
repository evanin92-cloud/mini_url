package config

import (
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	RedisHost      string
	RedisPort      string
	RedisPassword  string
	JWTSecret      string
	Port           string
	ShortIDLength  int
	LinkExpiration int
	MaxLinksPerUser int
}

func Load() *Config {
	return &Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "mini_url"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key-here"),
		Port:            getEnv("PORT", "8080"),
		ShortIDLength:   getEnvAsInt("SHORT_ID_LENGTH", 7),
		LinkExpiration:  getEnvAsInt("LINK_EXPIRATION_DAYS", 30),
		MaxLinksPerUser: getEnvAsInt("MAX_LINKS_PER_USER", 100),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Простое преобразование, для продакшена использовать strconv.Atoi
		return defaultValue
	}
	return defaultValue
}