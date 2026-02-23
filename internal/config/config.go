package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort        string
	SensitiveWordURL  string
	LLMServiceURL     string
	RequestTimeout    time.Duration
}

func Load() *Config {
	timeout := 30
	if t := os.Getenv("REQUEST_TIMEOUT"); t != "" {
		if v, err := strconv.Atoi(t); err == nil {
			timeout = v
		}
	}

	return &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		SensitiveWordURL: getEnv("SENSITIVE_WORD_URL", ""),
		LLMServiceURL:    getEnv("LLM_SERVICE_URL", ""),
		RequestTimeout:   time.Duration(timeout) * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
