package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		env     map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "默认配置",
			setup: func(t *testing.T) string {
				return ""
			},
			env: map[string]string{
				"SENSITIVE_WORD_URL": "http://sw.local",
				"LLM_SERVICE_URL":    "http://llm.local",
			},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "http://sw.local",
				LLMServiceURL:    "http://llm.local",
				RequestTimeout:   30 * time.Second,
				LogDir:           "./logs",
				LogMaxBackups:    30,
			},
		},
		{
			name: "环境变量覆盖",
			setup: func(t *testing.T) string {
				return ""
			},
			env: map[string]string{
				"SERVER_PORT":        "9090",
				"SENSITIVE_WORD_URL": "http://sw.local",
				"LLM_SERVICE_URL":    "http://llm.local",
				"REQUEST_TIMEOUT":    "60s",
			},
			want: &Config{
				ServerPort:       "9090",
				SensitiveWordURL: "http://sw.local",
				LLMServiceURL:    "http://llm.local",
				RequestTimeout:   60 * time.Second,
				LogDir:           "./logs",
				LogMaxBackups:    30,
			},
		},
		{
			name: "YAML配置",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				cfgPath := filepath.Join(tmpDir, "config.yaml")
				content := `log_dir: /var/log/app
log_max_backups: 10
`
				if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return cfgPath
			},
			env: map[string]string{
				"SENSITIVE_WORD_URL": "http://sw.local",
				"LLM_SERVICE_URL":    "http://llm.local",
			},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "http://sw.local",
				LLMServiceURL:    "http://llm.local",
				RequestTimeout:   30 * time.Second,
				LogDir:           "/var/log/app",
				LogMaxBackups:    10,
			},
		},
		{
			name: "REQUEST_TIMEOUT兼容性-duration格式",
			setup: func(t *testing.T) string {
				return ""
			},
			env: map[string]string{
				"SENSITIVE_WORD_URL": "http://sw.local",
				"LLM_SERVICE_URL":    "http://llm.local",
				"REQUEST_TIMEOUT":    "45s",
			},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "http://sw.local",
				LLMServiceURL:    "http://llm.local",
				RequestTimeout:   45 * time.Second,
				LogDir:           "./logs",
				LogMaxBackups:    30,
			},
		},
		{
			name: "REQUEST_TIMEOUT兼容性-纯数字格式",
			setup: func(t *testing.T) string {
				return ""
			},
			env: map[string]string{
				"SENSITIVE_WORD_URL": "http://sw.local",
				"LLM_SERVICE_URL":    "http://llm.local",
				"REQUEST_TIMEOUT":    "50",
			},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "http://sw.local",
				LLMServiceURL:    "http://llm.local",
				RequestTimeout:   50 * time.Second,
				LogDir:           "./logs",
				LogMaxBackups:    30,
			},
		},
		{
			name: "缺少必需字段-SENSITIVE_WORD_URL",
			setup: func(t *testing.T) string {
				return ""
			},
			env: map[string]string{
				"LLM_SERVICE_URL": "http://llm.local",
			},
			wantErr: true,
		},
		{
			name: "缺少必需字段-LLM_SERVICE_URL",
			setup: func(t *testing.T) string {
				return ""
			},
			env: map[string]string{
				"SENSITIVE_WORD_URL": "http://sw.local",
			},
			wantErr: true,
		},
		{
			name: "配置文件不存在-使用默认值",
			setup: func(t *testing.T) string {
				return "/nonexistent/config.yaml"
			},
			env: map[string]string{
				"SENSITIVE_WORD_URL": "http://sw.local",
				"LLM_SERVICE_URL":    "http://llm.local",
			},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "http://sw.local",
				LLMServiceURL:    "http://llm.local",
				RequestTimeout:   30 * time.Second,
				LogDir:           "./logs",
				LogMaxBackups:    30,
			},
		},
		{
			name: "YAML格式错误",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				cfgPath := filepath.Join(tmpDir, "config.yaml")
				content := `invalid: yaml: content: [[[`
				if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return cfgPath
			},
			env: map[string]string{
				"SENSITIVE_WORD_URL": "http://sw.local",
				"LLM_SERVICE_URL":    "http://llm.local",
			},
			wantErr: true,
		},
		{
			name: "CONFIG_PATH环境变量",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				cfgPath := filepath.Join(tmpDir, "custom.yaml")
				content := `log_dir: /custom/log
log_max_backups: 5
`
				if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return cfgPath
			},
			env: map[string]string{
				"SENSITIVE_WORD_URL": "http://sw.local",
				"LLM_SERVICE_URL":    "http://llm.local",
			},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "http://sw.local",
				LLMServiceURL:    "http://llm.local",
				RequestTimeout:   30 * time.Second,
				LogDir:           "/custom/log",
				LogMaxBackups:    5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			cfgPath := tt.setup(t)
			if cfgPath != "" {
				os.Setenv("CONFIG_PATH", cfgPath)
			}

			got, err := Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

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
