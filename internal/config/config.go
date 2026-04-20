package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO: consider splitting configs
type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBConnStr      string
	HTTPPort       string
	AllowedOrigins []string

	// Open AI configs for AI service
	OpenAIApiKey string

	// Input size constraints
	MaxUploadSizeBytes int64

	// Oauth
	OauthClientId     string
	OauthClientSecret string
}

func Load() *Config {
	config := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("POSTGRES_USER", "postgres"),
		DBPassword:     getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:         getEnv("POSTGRES_DB", "spend_wise"),
		HTTPPort:       getEnv("HTTP_PORT", "8000"),
		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:5173"), ","),

		// OpenAI configs for AI service
		OpenAIApiKey: os.Getenv("OPEN_AI_API_KEY"),

		// Oauth
		OauthClientId:     os.Getenv("OAUTH_CLIENT_ID"),
		OauthClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	}

	// Input size constraints
	maxUploadSizeKB, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_KB", "20"), 10, 64)
	if err != nil {
		maxUploadSizeKB = 5
	}
	config.MaxUploadSizeBytes = maxUploadSizeKB * 1024

	// DB connection url
	config.DBConnStr = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName,
	)

	return config
}

// getEnv reads an env var with a fallback default
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
