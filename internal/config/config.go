package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBConnStr  string
	HTTPPort   string

	// Open AI configs (AI)
	OpenAIUrl       string
	OpenAIModel     string
	OpenAIVersion   string
	OpenAIApiKey    string
	OpenAIMaxTokens int
}

func Load() *Config {
	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "postgres"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:     getEnv("POSTGRES_DB", "spend_wise"),
		HTTPPort:   getEnv("HTTP_PORT", "8000"),

		// OpenAI configs for AI
		OpenAIUrl:     getEnv("OPEN_AI_URL", "https://api.anthropic.com/v1/messages"),
		OpenAIModel:   getEnv("OPEN_AI_MODEL", "claude-sonnet-4-6"),
		OpenAIVersion: getEnv("OPEN_AI_VERSION", "2023-06-01"),
		OpenAIApiKey:  os.Getenv("OPEN_AI_API_KEY"),
	}

	config.DBConnStr = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName,
	)

	// Parse 'OPEN_AI_MAX_TOKENS' to int
	maxTokens, err := strconv.Atoi(getEnv("OPEN_AI_MAX_TOKENS", "1024"))

	if err != nil {
		maxTokens = 1024
	}

	config.OpenAIMaxTokens = maxTokens

	return config
}

// getEnv reads an env var with a fallback default
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
