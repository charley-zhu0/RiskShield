package service

import (
	"context"
	"errors"
	"testing"
)

func TestSensitiveWordDetector(t *testing.T) {
	tests := []struct {
		name       string
		hit        bool
		words      []string
		err        error
		wantReject bool
		wantReason string
		wantErr    bool
	}{
		{
			name:       "命中敏感词",
			hit:        true,
			words:      []string{"敏感词1", "敏感词2"},
			wantReject: true,
			wantReason: "命中敏感词",
		},
		{
			name:       "未命中敏感词",
			hit:        false,
			words:      []string{},
			wantReject: false,
		},
		{
			name:    "检测失败",
			err:     errors.New("service error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSensitiveWordClient{hit: tt.hit, words: tt.words, err: tt.err}
			detector := NewSensitiveWordDetector(mock)

			result, err := detector.Detect(context.Background(), "test content")

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.ShouldReject != tt.wantReject {
				t.Errorf("ShouldReject = %v; want %v", result.ShouldReject, tt.wantReject)
			}

			if tt.wantReject && result.Reason != tt.wantReason {
				t.Errorf("Reason = %q; want %q", result.Reason, tt.wantReason)
			}
		})
	}
}

func TestLLMDetector(t *testing.T) {
	tests := []struct {
		name       string
		safetyType string
		attackType string
		err        error
		wantReject bool
		wantErr    bool
	}{
		{
			name:       "正常内容",
			safetyType: "正常",
			attackType: "无",
			wantReject: false,
		},
		{
			name:       "违规内容",
			safetyType: "色情",
			attackType: "无",
			wantReject: true,
		},
		{
			name:    "分类失败",
			err:     errors.New("llm error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockLLMClient{
				safetyType: tt.safetyType,
				attackType: tt.attackType,
				err:        tt.err,
			}
			detector := NewLLMDetector(mock)

			result, err := detector.Detect(context.Background(), "test content")

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.ShouldReject != tt.wantReject {
				t.Errorf("ShouldReject = %v; want %v", result.ShouldReject, tt.wantReject)
			}
		})
	}
}
