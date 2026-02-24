package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort       string
	SensitiveWordURL string
	LLMServiceURL    string
	RequestTimeout   time.Duration
	LogDir           string
	LogMaxBackups    int
}

func Load() *Config {
	timeout := 30
	if t := os.Getenv("REQUEST_TIMEOUT"); t != "" {
		if v, err := strconv.Atoi(t); err == nil {
			timeout = v
		}
	}

	logMaxBackups := 30
	if b := os.Getenv("LOG_MAX_BACKUPS"); b != "" {
		if v, err := strconv.Atoi(b); err == nil {
			logMaxBackups = v
		}
	}

	return &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		SensitiveWordURL: getEnv("SENSITIVE_WORD_URL", ""),
		LLMServiceURL:    getEnv("LLM_SERVICE_URL", ""),
		RequestTimeout:   time.Duration(timeout) * time.Second,
		LogDir:           getEnv("LOG_DIR", "logs"),
		LogMaxBackups:    logMaxBackups,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
