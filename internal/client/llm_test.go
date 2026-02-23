package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLLMClient_Classify(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		serverResp     string
		statusCode     int
		wantSafetyType string
		wantAttackType string
		wantErr        bool
	}{
		{
			name:           "检测到违规内容",
			content:        "1989年暴动事件",
			serverResp:     `{"choices":[{"message":{"content":"{\"safety_type\":\"涉政:敏感政治事件:六四事件\",\"attack_type\":\"无\"}"}}]}`,
			statusCode:     http.StatusOK,
			wantSafetyType: "涉政:敏感政治事件:六四事件",
			wantAttackType: "无",
		},
		{
			name:           "正常内容",
			content:        "今天天气不错",
			serverResp:     `{"choices":[{"message":{"content":"{\"safety_type\":\"正常\",\"attack_type\":\"无\"}"}}]}`,
			statusCode:     http.StatusOK,
			wantSafetyType: "正常",
			wantAttackType: "无",
		},
		{
			name:       "服务返回错误",
			content:    "测试",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "无效JSON响应",
			content:    "测试",
			serverResp: `{invalid}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "空choices",
			content:    "测试",
			serverResp: `{"choices":[]}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("期望 POST 请求, 得到 %s", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				if tt.serverResp != "" {
					w.Write([]byte(tt.serverResp))
				}
			}))
			defer server.Close()

			client := NewLLMClient(server.URL, 5*time.Second)
			safetyType, attackType, err := client.Classify(context.Background(), tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("Classify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if safetyType != tt.wantSafetyType {
					t.Errorf("Classify() safetyType = %v, want %v", safetyType, tt.wantSafetyType)
				}
				if attackType != tt.wantAttackType {
					t.Errorf("Classify() attackType = %v, want %v", attackType, tt.wantAttackType)
				}
			}
		})
	}
}
