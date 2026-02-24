package service

import (
	"context"
	"time"

	"github.com/charley/riskshield/internal/client"
	"github.com/charley/riskshield/internal/domain"
)

type ModerationService interface {
	Moderate(ctx context.Context, content string) (*domain.ModerationResult, error)
}

type moderationService struct {
	executor   *ConcurrentExecutor
	aggregator *Aggregator
}

func NewModerationService(swClient client.SensitiveWordClient, llmClient client.LLMClient) ModerationService {
	detectors := []domain.Detector{
		NewSensitiveWordDetector(swClient),
		NewLLMDetector(llmClient),
	}

	return &moderationService{
		executor:   NewConcurrentExecutor(detectors, 300*time.Millisecond),
		aggregator: NewAggregator(),
	}
}

func (s *moderationService) Moderate(ctx context.Context, content string) (*domain.ModerationResult, error) {
	results := s.executor.Execute(ctx, content)
	return s.aggregator.Aggregate(results), nil
}
