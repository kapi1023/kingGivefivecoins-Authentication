package config

import (
	"log"
	"os"
)

type Config struct {
	Port           string
	JWTKey         string
	GoogleClient   string
	GoogleSecret   string
	LinkedInClient string
	LinkedInSecret string
	AppleClient    string
	AppleSecret    string
	DatabaseURL    string
}

func Load() *Config {
	getEnv := func(key string, optional ...bool) string {
		value := os.Getenv(key)
		if value == "" && len(optional) == 0 {
			log.Fatalf("Environment variable %s not set", key)
		}
		return value
	}

	config := &Config{
		Port:           getEnv("PORT"),
		JWTKey:         getEnv("JWT_KEY"),
		GoogleClient:   getEnv("GOOGLE_CLIENT_ID"),
		GoogleSecret:   getEnv("GOOGLE_CLIENT_SECRET"),
		LinkedInClient: getEnv("LINKEDIN_CLIENT_ID"),
		LinkedInSecret: getEnv("LINKEDIN_CLIENT_SECRET"),
		AppleClient:    getEnv("APPLE_CLIENT_ID"),
		AppleSecret:    getEnv("APPLE_CLIENT_SECRET"),
		DatabaseURL:    getEnv("DATABASE_URL"),
	}

	if config.GoogleClient == "" || config.GoogleSecret == "" {
		log.Fatal("Google OAuth2 credentials are not set")
	}
	return config
}
