package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLoadLogConfigFromJSON tests loading log configuration from JSON file
func TestLoadLogConfigFromJSON(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		jsonData   string
		want       *LogConfig
		wantErr    bool
	}{
		{
			name:       "valid config with all fields",
			configPath: "test_valid.json",
			jsonData: `{
				"log": {
					"dir": "/custom/logs",
					"max_backups": 50
				}
			}`,
			want: &LogConfig{
				Dir:        "/custom/logs",
				MaxBackups: 50,
			},
			wantErr: false,
		},
		{
			name:       "valid config with default values",
			configPath: "test_defaults.json",
			jsonData: `{
				"log": {}
			}`,
			want: &LogConfig{
				Dir:        "./logs",
				MaxBackups: 30,
			},
			wantErr: false,
		},
		{
			name:       "config file does not exist",
			configPath: "nonexistent.json",
			jsonData:   "",
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "invalid JSON format",
			configPath: "test_invalid.json",
			jsonData:   `{invalid json}`,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "empty file",
			configPath: "test_empty.json",
			jsonData:   "",
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "missing log section",
			configPath: "test_no_log.json",
			jsonData:   `{"other": "data"}`,
			want: &LogConfig{
				Dir:        "./logs",
				MaxBackups: 30,
			},
			wantErr: false,
		},
		{
			name:       "partial config - only dir",
			configPath: "test_partial_dir.json",
			jsonData: `{
				"log": {
					"dir": "/tmp/logs"
				}
			}`,
			want: &LogConfig{
				Dir:        "/tmp/logs",
				MaxBackups: 30,
			},
			wantErr: false,
		},
		{
			name:       "partial config - only max_backups",
			configPath: "test_partial_backups.json",
			jsonData: `{
				"log": {
					"max_backups": 100
				}
			}`,
			want: &LogConfig{
				Dir:        "./logs",
				MaxBackups: 100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test files
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, tt.configPath)

			// Create test file if jsonData is provided
			if tt.jsonData != "" {
				if err := os.WriteFile(configPath, []byte(tt.jsonData), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			got, err := LoadLogConfig(configPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadLogConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil && tt.want != nil {
				if got.Dir != tt.want.Dir {
					t.Errorf("LogConfig.Dir = %v, want %v", got.Dir, tt.want.Dir)
				}
				if got.MaxBackups != tt.want.MaxBackups {
					t.Errorf("LogConfig.MaxBackups = %v, want %v", got.MaxBackups, tt.want.MaxBackups)
				}
			}
		})
	}
}

// TestLoad_WithJSONConfig tests the integrated Load function with JSON config
func TestLoad_WithJSONConfig(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		envVars  map[string]string
		want     *Config
		wantErr  bool
	}{
		{
			name: "load with JSON config and env vars",
			jsonData: `{
				"log": {
					"dir": "./logs",
					"max_backups": 30
				}
			}`,
			envVars: map[string]string{
				"SERVER_PORT":        "9090",
				"SENSITIVE_WORD_URL": "http://sw.example.com",
				"LLM_SERVICE_URL":    "http://llm.example.com",
				"REQUEST_TIMEOUT":    "60",
			},
			want: &Config{
				ServerPort:       "9090",
				SensitiveWordURL: "http://sw.example.com",
				LLMServiceURL:    "http://llm.example.com",
				RequestTimeout:   60 * time.Second,
				LogDir:           "./logs",
				LogMaxBackups:    30,
			},
			wantErr: false,
		},
		{
			name: "load with default JSON config",
			jsonData: `{
				"log": {
					"dir": "./logs",
					"max_backups": 30
				}
			}`,
			envVars: map[string]string{},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "",
				LLMServiceURL:    "",
				RequestTimeout:   30 * time.Second,
				LogDir:           "./logs",
				LogMaxBackups:    30,
			},
			wantErr: false,
		},
		{
			name: "load with custom log dir",
			jsonData: `{
				"log": {
					"dir": "/var/log/riskshield",
					"max_backups": 60
				}
			}`,
			envVars: map[string]string{
				"SERVER_PORT": "8080",
			},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "",
				LLMServiceURL:    "",
				RequestTimeout:   30 * time.Second,
				LogDir:           "/var/log/riskshield",
				LogMaxBackups:    60,
			},
			wantErr: false,
		},
		{
			name:     "fail when config file missing",
			jsonData: "",
			envVars:  map[string]string{},
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Clearenv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Create temporary config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "risk.json")

			if tt.jsonData != "" {
				if err := os.WriteFile(configPath, []byte(tt.jsonData), 0644); err != nil {
					t.Fatalf("failed to create test config file: %v", err)
				}
			}

			// Override config path for testing
			originalPath := DefaultConfigPath
			DefaultConfigPath = configPath
			defer func() {
				DefaultConfigPath = originalPath
			}()

			got, err := LoadWithPath(configPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil && tt.want != nil {
				if got.ServerPort != tt.want.ServerPort {
					t.Errorf("Config.ServerPort = %v, want %v", got.ServerPort, tt.want.ServerPort)
				}
				if got.SensitiveWordURL != tt.want.SensitiveWordURL {
					t.Errorf("Config.SensitiveWordURL = %v, want %v", got.SensitiveWordURL, tt.want.SensitiveWordURL)
				}
				if got.LLMServiceURL != tt.want.LLMServiceURL {
					t.Errorf("Config.LLMServiceURL = %v, want %v", got.LLMServiceURL, tt.want.LLMServiceURL)
				}
				if got.RequestTimeout != tt.want.RequestTimeout {
					t.Errorf("Config.RequestTimeout = %v, want %v", got.RequestTimeout, tt.want.RequestTimeout)
				}
				if got.LogDir != tt.want.LogDir {
					t.Errorf("Config.LogDir = %v, want %v", got.LogDir, tt.want.LogDir)
				}
				if got.LogMaxBackups != tt.want.LogMaxBackups {
					t.Errorf("Config.LogMaxBackups = %v, want %v", got.LogMaxBackups, tt.want.LogMaxBackups)
				}
			}
		})
	}
}

// TestLoad_LegacyBehavior ensures backward compatibility - SHOULD BE REMOVED
// This test documents the old behavior but will be removed after migration
func TestLoad_LegacyBehavior(t *testing.T) {
	t.Skip("Legacy test - will be removed after migration to JSON config")

	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
	}{
		{
			name:    "使用默认值",
			envVars: map[string]string{},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "",
				LLMServiceURL:    "",
				RequestTimeout:   30 * time.Second,
				LogDir:           "logs",
				LogMaxBackups:    30,
			},
		},
		{
			name: "使用环境变量",
			envVars: map[string]string{
				"SERVER_PORT":        "9090",
				"SENSITIVE_WORD_URL": "http://sw.example.com",
				"LLM_SERVICE_URL":    "http://llm.example.com",
				"REQUEST_TIMEOUT":    "60",
				"LOG_DIR":            "/var/log/app",
				"LOG_MAX_BACKUPS":    "60",
			},
			want: &Config{
				ServerPort:       "9090",
				SensitiveWordURL: "http://sw.example.com",
				LLMServiceURL:    "http://llm.example.com",
				RequestTimeout:   60 * time.Second,
				LogDir:           "/var/log/app",
				LogMaxBackups:    60,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理环境变量
			os.Clearenv()

			// 设置测试环境变量
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			got := Load()

			if got.ServerPort != tt.want.ServerPort {
				t.Errorf("ServerPort = %v, want %v", got.ServerPort, tt.want.ServerPort)
			}
			if got.SensitiveWordURL != tt.want.SensitiveWordURL {
				t.Errorf("SensitiveWordURL = %v, want %v", got.SensitiveWordURL, tt.want.SensitiveWordURL)
			}
			if got.LLMServiceURL != tt.want.LLMServiceURL {
				t.Errorf("LLMServiceURL = %v, want %v", got.LLMServiceURL, tt.want.LLMServiceURL)
			}
			if got.RequestTimeout != tt.want.RequestTimeout {
				t.Errorf("RequestTimeout = %v, want %v", got.RequestTimeout, tt.want.RequestTimeout)
			}
			if got.LogDir != tt.want.LogDir {
				t.Errorf("LogDir = %v, want %v", got.LogDir, tt.want.LogDir)
			}
			if got.LogMaxBackups != tt.want.LogMaxBackups {
				t.Errorf("LogMaxBackups = %v, want %v", got.LogMaxBackups, tt.want.LogMaxBackups)
			}
		})
	}
}

// TestLoad tests the default Load function with risk.json
func TestLoad(t *testing.T) {
	// Save original DefaultConfigPath
	originalPath := DefaultConfigPath
	defer func() {
		DefaultConfigPath = originalPath
	}()

	t.Run("successfully loads from default config", func(t *testing.T) {
		// Create temporary config
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "risk.json")
		jsonData := `{
			"log": {
				"dir": "./logs",
				"max_backups": 30
			}
		}`
		if err := os.WriteFile(configPath, []byte(jsonData), 0644); err != nil {
			t.Fatalf("failed to create test config: %v", err)
		}

		// Override default path
		DefaultConfigPath = configPath

		// Clear env
		os.Clearenv()

		// Should not panic
		cfg := Load()

		if cfg == nil {
			t.Fatal("expected config, got nil")
		}

		if cfg.LogDir != "./logs" {
			t.Errorf("LogDir = %v, want ./logs", cfg.LogDir)
		}
		if cfg.LogMaxBackups != 30 {
			t.Errorf("LogMaxBackups = %v, want 30", cfg.LogMaxBackups)
		}
	})

	t.Run("panics when config file is missing", func(t *testing.T) {
		// Set non-existent path
		DefaultConfigPath = "/nonexistent/path/risk.json"

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()

		Load()
	})

	t.Run("panics when config file is invalid", func(t *testing.T) {
		// Create invalid config
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "risk.json")
		if err := os.WriteFile(configPath, []byte("{invalid}"), 0644); err != nil {
			t.Fatalf("failed to create test config: %v", err)
		}

		DefaultConfigPath = configPath

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()

		Load()
	})
}
