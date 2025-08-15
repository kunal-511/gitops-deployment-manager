package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string `json:"port"`
	DatabaseURL string `json:"database_url"`
	Env         string `json:"env"`
	JWTSecret   string `json:"jwt"`
}

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://neondb_owner:npg_NPDLs3H4mFqE@ep-sparkling-cherry-ad7zimxj-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"),
		Env:         getEnv("ENV", "dev"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
	}
	log.Printf("Config loaded: env=%s port=%s", cfg.Env, cfg.Port)
	return cfg
}

func getEnv(k, def string) string {
	if value := os.Getenv(k); value != "" {
		return value
	}
	return def
}
