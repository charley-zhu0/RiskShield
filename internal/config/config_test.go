package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
	}{
		{
			name: "使用默认值",
			envVars: map[string]string{},
			want: &Config{
				ServerPort:       "8080",
				SensitiveWordURL: "",
				LLMServiceURL:    "",
				RequestTimeout:   30 * time.Second,
			},
		},
		{
			name: "使用环境变量",
			envVars: map[string]string{
				"SERVER_PORT":         "9090",
				"SENSITIVE_WORD_URL":  "http://sw.example.com",
				"LLM_SERVICE_URL":     "http://llm.example.com",
				"REQUEST_TIMEOUT":     "60",
			},
			want: &Config{
				ServerPort:       "9090",
				SensitiveWordURL: "http://sw.example.com",
				LLMServiceURL:    "http://llm.example.com",
				RequestTimeout:   60 * time.Second,
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
		})
	}
}
