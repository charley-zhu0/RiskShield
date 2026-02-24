package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// DefaultConfigPath is the default path to the risk.json configuration file
var DefaultConfigPath = "./internal/config/risk.json"

// Config holds all application configuration
type Config struct {
	ServerPort       string
	SensitiveWordURL string
	LLMServiceURL    string
	RequestTimeout   time.Duration
	LogDir           string
	LogMaxBackups    int
}

// LogConfig holds log-specific configuration from JSON file
type LogConfig struct {
	Dir        string `json:"dir"`
	MaxBackups int    `json:"max_backups"`
}

// RiskConfig represents the structure of risk.json file
type RiskConfig struct {
	Log LogConfig `json:"log"`
}

// LoadLogConfig loads log configuration from JSON file
func LoadLogConfig(configPath string) (*LogConfig, error) {
	// Read the JSON file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Check for empty file
	if len(data) == 0 {
		return nil, fmt.Errorf("config file %s is empty", configPath)
	}

	// Parse JSON
	var riskConfig RiskConfig
	if err := json.Unmarshal(data, &riskConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Apply default values if not specified
	logConfig := &LogConfig{
		Dir:        riskConfig.Log.Dir,
		MaxBackups: riskConfig.Log.MaxBackups,
	}

	if logConfig.Dir == "" {
		logConfig.Dir = "./logs"
	}

	if logConfig.MaxBackups == 0 {
		logConfig.MaxBackups = 30
	}

	return logConfig, nil
}

// Load loads configuration from environment variables and JSON file
// This function loads log configuration from risk.json and other configs from env vars
func Load() *Config {
	configPath := DefaultConfigPath
	if absPath, err := filepath.Abs(configPath); err == nil {
		configPath = absPath
	}

	cfg, err := LoadWithPath(configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	return cfg
}

// LoadWithPath loads configuration from a specific config file path
func LoadWithPath(configPath string) (*Config, error) {
	// Load log configuration from JSON file
	logConfig, err := LoadLogConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load log config: %w", err)
	}

	// Load other configuration from environment variables
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
		LogDir:           logConfig.Dir,
		LogMaxBackups:    logConfig.MaxBackups,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
