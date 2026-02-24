package service

import (
	"context"

	"github.com/charley/riskshield/internal/client"
	"github.com/charley/riskshield/internal/domain"
)

type sensitiveWordDetector struct {
	client client.SensitiveWordClient
}

func NewSensitiveWordDetector(client client.SensitiveWordClient) domain.Detector {
	return &sensitiveWordDetector{client: client}
}

func (d *sensitiveWordDetector) Detect(ctx context.Context, content string) (*domain.DetectionResult, error) {
	hit, words, err := d.client.Match(ctx, content)
	if err != nil {
		return nil, err
	}

	if hit {
		return &domain.DetectionResult{
			ShouldReject: true,
			Reason:       "命中敏感词",
			HitWords:     words,
		}, nil
	}

	return &domain.DetectionResult{
		ShouldReject: false,
		HitWords:     []string{},
	}, nil
}

func (d *sensitiveWordDetector) Name() string {
	return "SensitiveWord"
}

type llmDetector struct {
	client client.LLMClient
}

func NewLLMDetector(client client.LLMClient) domain.Detector {
	return &llmDetector{client: client}
}

func (d *llmDetector) Detect(ctx context.Context, content string) (*domain.DetectionResult, error) {
	safetyType, attackType, err := d.client.Classify(ctx, content)
	if err != nil {
		return nil, err
	}

	shouldReject := safetyType != "正常"
	return &domain.DetectionResult{
		ShouldReject: shouldReject,
		Reason:       "内容违规",
		HitWords:     []string{},
		SafetyType:   safetyType,
		AttackType:   attackType,
	}, nil
}

func (d *llmDetector) Name() string {
	return "LLM"
}
