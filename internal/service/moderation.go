package service

import (
	"context"

	"github.com/charley/riskshield/internal/client"
	"github.com/charley/riskshield/internal/domain"
)

type ModerationService interface {
	Moderate(ctx context.Context, content string) (*domain.ModerationResult, error)
}

type moderationService struct {
	swClient  client.SensitiveWordClient
	llmClient client.LLMClient
}

func NewModerationService(swClient client.SensitiveWordClient, llmClient client.LLMClient) ModerationService {
	return &moderationService{
		swClient:  swClient,
		llmClient: llmClient,
	}
}

func (s *moderationService) Moderate(ctx context.Context, content string) (*domain.ModerationResult, error) {
	// 检查敏感词
	hit, words, err := s.swClient.Match(ctx, content)
	if err != nil {
		return nil, err
	}

	if hit {
		return &domain.ModerationResult{
			Decision: "REJECT",
			Reason:   "命中敏感词",
			HitWords: words,
		}, nil
	}

	// 大模型分类
	safetyType, attackType, err := s.llmClient.Classify(ctx, content)
	if err != nil {
		return nil, err
	}

	result := &domain.ModerationResult{
		HitWords:   []string{},
		SafetyType: safetyType,
		AttackType: attackType,
	}

	if safetyType != "正常" {
		result.Decision = "REJECT"
		result.Reason = "内容违规"
	} else {
		result.Decision = "PASS"
	}

	return result, nil
}
