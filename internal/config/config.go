package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
	OpenAIUrl     string
	OpenAIModel   string
	OpenAIVersion string
	OpenAIApiKey  string

	// Input size constraints
	MaxUploadSizeBytes int64

	// Oauth
	OauthApiUrl        string
	AuthSessionLifeSec time.Duration
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
		OpenAIUrl:     getEnv("OPEN_AI_URL", "https://api.anthropic.com/v1/messages"),
		OpenAIModel:   getEnv("OPEN_AI_MODEL", "claude-sonnet-4-6"),
		OpenAIVersion: getEnv("OPEN_AI_VERSION", "2023-06-01"),
		OpenAIApiKey:  os.Getenv("OPEN_AI_API_KEY"),

		// Oauth
		OauthApiUrl:        getEnv("OAUTH_API_URL", "https://www.googleapis.com/oauth2/v2"),
		AuthSessionLifeSec: 7 * 24 * time.Hour, // 7 days
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
