package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort       string        `mapstructure:"server_port"`
	SensitiveWordURL string        `mapstructure:"sensitive_word_url"`
	LLMServiceURL    string        `mapstructure:"llm_service_url"`
	RequestTimeout   time.Duration `mapstructure:"request_timeout"`
	LogDir           string        `mapstructure:"log_dir"`
	LogMaxBackups    int           `mapstructure:"log_max_backups"`
	MySQLDSN         string        `mapstructure:"mysql_dsn"`
}

func Load() (*Config, error) {
	v := viper.New()
	setDefaults(v)
	bindEnvVars(v)

	configPath := v.GetString("config_path")
	if configPath != "" {
		v.SetConfigFile(configPath)
		err := v.ReadInConfig()
		if err != nil {
			// Check if it's a "file not found" error using string matching
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// File not found is OK
			} else if os.IsNotExist(err) {
				// Also handle os.PathError
			} else {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	// Handle REQUEST_TIMEOUT compatibility before unmarshal
	if timeoutStr := v.GetString("request_timeout"); timeoutStr != "" {
		if seconds, err := strconv.Atoi(timeoutStr); err == nil {
			v.Set("request_timeout", time.Duration(seconds)*time.Second)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server_port", "8080")
	v.SetDefault("request_timeout", 30*time.Second)
	v.SetDefault("log_dir", "./logs")
	v.SetDefault("log_max_backups", 30)
}

func bindEnvVars(v *viper.Viper) {
	v.BindEnv("config_path", "CONFIG_PATH")
	v.BindEnv("server_port", "SERVER_PORT")
	v.BindEnv("sensitive_word_url", "SENSITIVE_WORD_URL")
	v.BindEnv("llm_service_url", "LLM_SERVICE_URL")
	v.BindEnv("request_timeout", "REQUEST_TIMEOUT")
	v.BindEnv("log_dir", "LOG_DIR")
	v.BindEnv("log_max_backups", "LOG_MAX_BACKUPS")
	v.BindEnv("mysql_dsn", "MYSQL_DSN")
}

func validate(cfg *Config) error {
	if cfg.SensitiveWordURL == "" {
		return fmt.Errorf("SENSITIVE_WORD_URL is required")
	}
	if cfg.LLMServiceURL == "" {
		return fmt.Errorf("LLM_SERVICE_URL is required")
	}
	if cfg.MySQLDSN == "" {
		return fmt.Errorf("MYSQL_DSN is required")
	}
	return nil
}
