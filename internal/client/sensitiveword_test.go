package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSensitiveWordClient_Match(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		serverResp string
		statusCode int
		wantHit    bool
		wantWords  []string
		wantErr    bool
	}{
		{
			name:       "命中敏感词",
			content:    "测试敏感词",
			serverResp: `{"query":"测试敏感词","hit_words":["敏感词"]}`,
			statusCode: http.StatusOK,
			wantHit:    true,
			wantWords:  []string{"敏感词"},
		},
		{
			name:       "未命中敏感词",
			content:    "正常内容",
			serverResp: `{"query":"正常内容","hit_words":[]}`,
			statusCode: http.StatusOK,
			wantHit:    false,
			wantWords:  []string{},
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
			serverResp: `{invalid json}`,
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
				if r.URL.Query().Get("tenantId") != "1" {
					t.Errorf("期望 tenantId=1")
				}
				w.WriteHeader(tt.statusCode)
				if tt.serverResp != "" {
					w.Write([]byte(tt.serverResp))
				}
			}))
			defer server.Close()

			client := NewSensitiveWordClient(server.URL, 5*time.Second)
			hit, words, err := client.Match(context.Background(), tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hit != tt.wantHit {
					t.Errorf("Match() hit = %v, want %v", hit, tt.wantHit)
				}
				if len(words) != len(tt.wantWords) {
					t.Errorf("Match() words = %v, want %v", words, tt.wantWords)
				}
			}
		})
	}
}
