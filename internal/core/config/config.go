package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv            string
	Port              string
	DatabaseURL       string
	JWTSecret         string
	JWTAccessTTLMin   int
	JWTRefreshTTLHour int
	CORSOrigins       string
	RateLimitStore    string
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
	LoginLimit        int
	LoginWindowMin    int
}

func Load() Config {
	return Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=football_api port=5432 sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
		JWTAccessTTLMin:   getEnvInt("JWT_ACCESS_TTL_MIN", 30),
		JWTRefreshTTLHour: getEnvInt("JWT_REFRESH_TTL_HOUR", 72),
		CORSOrigins:       getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:8081"),
		RateLimitStore:    getEnv("RATE_LIMIT_STORE", "memory"),
		RedisAddr:         getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           getEnvInt("REDIS_DB", 0),
		LoginLimit:        getEnvInt("LOGIN_LIMIT", 5),
		LoginWindowMin:    getEnvInt("LOGIN_WINDOW_MIN", 15),
	}
}

func getEnv(k, fallback string) string {
	v := os.Getenv(k)
	if v == "" {
		return fallback
	}
	return v
}

func getEnvInt(k string, fallback int) int {
	v := os.Getenv(k)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
