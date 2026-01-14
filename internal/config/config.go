package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	Port             string
	DatabaseURL      string
	CORSOrigins      []string
	JWTSecret        string
	JWTRefreshSecret string
	AccessTTL        time.Duration
	RefreshTTL       time.Duration
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsOrigins := os.Getenv("CORS_ORIGINS")
	var origins []string
	if corsOrigins != "" {
		origins = strings.Split(corsOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
	} else {
		origins = []string{"http://localhost:3000"}
	}

	accessTTL := parseDuration(os.Getenv("ACCESS_TTL"), 15*time.Minute)
	refreshTTL := parseDuration(os.Getenv("REFRESH_TTL"), 168*time.Hour) // 7 days

	return &Config{
		Port:             port,
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		CORSOrigins:      origins,
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		AccessTTL:        accessTTL,
		RefreshTTL:       refreshTTL,
	}
}

func parseDuration(s string, defaultDuration time.Duration) time.Duration {
	if s == "" {
		return defaultDuration
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultDuration
	}
	return d
}
