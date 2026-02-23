package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/charley/riskshield/internal/domain"
	"github.com/gin-gonic/gin"
)

type mockModerationService struct {
	result *domain.ModerationResult
	err    error
}

func (m *mockModerationService) Moderate(ctx context.Context, content string) (*domain.ModerationResult, error) {
	return m.result, m.err
}

func TestModerationHandler_Moderate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		formData       map[string]string
		mockResult     *domain.ModerationResult
		mockErr        error
		wantStatusCode int
		wantDecision   string
	}{
		{
			name: "成功-通过",
			formData: map[string]string{
				"model":   "gpt-4",
				"content": "今天天气不错",
			},
			mockResult: &domain.ModerationResult{
				Decision:   "PASS",
				SafetyType: "正常",
				AttackType: "无",
				HitWords:   []string{},
			},
			wantStatusCode: http.StatusOK,
			wantDecision:   "PASS",
		},
		{
			name: "成功-拒绝",
			formData: map[string]string{
				"model":   "gpt-4",
				"content": "敏感内容",
			},
			mockResult: &domain.ModerationResult{
				Decision: "REJECT",
				Reason:   "命中敏感词",
				HitWords: []string{"敏感"},
			},
			wantStatusCode: http.StatusOK,
			wantDecision:   "REJECT",
		},
		{
			name: "缺少必填参数",
			formData: map[string]string{
				"model": "gpt-4",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockModerationService{
				result: tt.mockResult,
				err:    tt.mockErr,
			}
			handler := NewModerationHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			form := url.Values{}
			for k, v := range tt.formData {
				form.Set(k, v)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/private/moderations/text", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			c.Request = req

			handler.Moderate(c)

			if w.Code != tt.wantStatusCode {
				t.Errorf("状态码 = %v, want %v", w.Code, tt.wantStatusCode)
			}
		})
	}
}
