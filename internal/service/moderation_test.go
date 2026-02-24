package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockSensitiveWordClient struct {
	hit   bool
	words []string
	err   error
	delay time.Duration
}

func (m *mockSensitiveWordClient) Match(ctx context.Context, content string) (bool, []string, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return m.hit, m.words, m.err
}

type mockLLMClient struct {
	safetyType string
	attackType string
	err        error
	delay      time.Duration
}

func (m *mockLLMClient) Classify(ctx context.Context, content string) (string, string, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return m.safetyType, m.attackType, m.err
}

func TestModerationService_Moderate(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		swHit        bool
		swWords      []string
		swErr        error
		llmSafety    string
		llmAttack    string
		llmErr       error
		wantDecision string
		wantReason   string
	}{
		{
			name:         "敏感词命中-拒绝",
			content:      "测试敏感词",
			swHit:        true,
			swWords:      []string{"敏感词"},
			llmSafety:    "正常",
			llmAttack:    "无",
			wantDecision: "REJECT",
			wantReason:   "命中敏感词",
		},
		{
			name:         "大模型判定违规-拒绝",
			content:      "1989年暴动",
			swHit:        false,
			swWords:      []string{},
			llmSafety:    "涉政:敏感政治事件:六四事件",
			llmAttack:    "无",
			wantDecision: "REJECT",
			wantReason:   "内容违规",
		},
		{
			name:         "两者都通过-放行",
			content:      "今天天气不错",
			swHit:        false,
			swWords:      []string{},
			llmSafety:    "正常",
			llmAttack:    "无",
			wantDecision: "PASS",
			wantReason:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			swClient := &mockSensitiveWordClient{
				hit:   tt.swHit,
				words: tt.swWords,
				err:   tt.swErr,
			}
			llmClient := &mockLLMClient{
				safetyType: tt.llmSafety,
				attackType: tt.llmAttack,
				err:        tt.llmErr,
			}

			svc := NewModerationService(swClient, llmClient)
			result, err := svc.Moderate(context.Background(), tt.content)

			if err != nil {
				t.Fatalf("Moderate() error = %v", err)
			}
			if result.Decision != tt.wantDecision {
				t.Errorf("Decision = %v, want %v", result.Decision, tt.wantDecision)
			}
			if result.Reason != tt.wantReason {
				t.Errorf("Reason = %v, want %v", result.Reason, tt.wantReason)
			}
		})
	}
}

func TestModerationService_Moderate_Concurrent(t *testing.T) {
	t.Run("并发执行快于串行", func(t *testing.T) {
		swClient := &mockSensitiveWordClient{
			hit:   false,
			words: []string{},
			delay: 100 * time.Millisecond,
		}
		llmClient := &mockLLMClient{
			safetyType: "正常",
			attackType: "无",
			delay:      100 * time.Millisecond,
		}

		svc := NewModerationService(swClient, llmClient)
		start := time.Now()
		_, err := svc.Moderate(context.Background(), "test")
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if elapsed > 200*time.Millisecond {
			t.Errorf("expected concurrent execution < 200ms, got %v", elapsed)
		}
	})

	t.Run("超时忽略慢服务", func(t *testing.T) {
		swClient := &mockSensitiveWordClient{
			hit:   false,
			words: []string{},
			delay: 50 * time.Millisecond,
		}
		llmClient := &mockLLMClient{
			safetyType: "正常",
			attackType: "无",
			delay:      500 * time.Millisecond,
		}

		svc := NewModerationService(swClient, llmClient)
		result, err := svc.Moderate(context.Background(), "test")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Decision != "PASS" {
			t.Errorf("expected PASS, got %s", result.Decision)
		}
	})
}

func TestModerationService_Moderate_Errors(t *testing.T) {
	t.Run("敏感词服务失败-继续执行", func(t *testing.T) {
		swClient := &mockSensitiveWordClient{err: errors.New("服务错误")}
		llmClient := &mockLLMClient{safetyType: "正常", attackType: "无"}
		svc := NewModerationService(swClient, llmClient)

		result, err := svc.Moderate(context.Background(), "测试")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Decision != "PASS" {
			t.Errorf("expected PASS, got %s", result.Decision)
		}
	})

	t.Run("大模型服务失败-继续执行", func(t *testing.T) {
		swClient := &mockSensitiveWordClient{hit: false, words: []string{}}
		llmClient := &mockLLMClient{err: errors.New("服务错误")}
		svc := NewModerationService(swClient, llmClient)

		result, err := svc.Moderate(context.Background(), "测试")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Decision != "PASS" {
			t.Errorf("expected PASS, got %s", result.Decision)
		}
	})
}
