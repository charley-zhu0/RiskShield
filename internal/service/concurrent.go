package service

import (
	"context"
	"time"

	"github.com/charley/riskshield/internal/domain"
)

type DetectionResult struct {
	DetectorName string
	Result       *domain.DetectionResult
	Error        error
}

type ConcurrentExecutor struct {
	detectors []domain.Detector
	timeout   time.Duration
}

func NewConcurrentExecutor(detectors []domain.Detector, timeout time.Duration) *ConcurrentExecutor {
	return &ConcurrentExecutor{
		detectors: detectors,
		timeout:   timeout,
	}
}

func (e *ConcurrentExecutor) Execute(ctx context.Context, content string) []DetectionResult {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	resultCh := make(chan DetectionResult, len(e.detectors))

	for _, detector := range e.detectors {
		go func(d domain.Detector) {
			result, err := d.Detect(ctx, content)
			select {
			case resultCh <- DetectionResult{
				DetectorName: d.Name(),
				Result:       result,
				Error:        err,
			}:
			case <-ctx.Done():
			}
		}(detector)
	}

	var results []DetectionResult
	for i := 0; i < len(e.detectors); i++ {
		select {
		case r := <-resultCh:
			results = append(results, r)
		case <-ctx.Done():
			return results
		}
	}

	return results
}

type Aggregator struct{}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Aggregate(results []DetectionResult) *domain.ModerationResult {
	final := &domain.ModerationResult{
		Decision: "PASS",
		HitWords: []string{},
	}

	for _, r := range results {
		if r.Error != nil {
			continue
		}

		if r.Result.ShouldReject {
			final.Decision = "REJECT"
			final.Reason = r.Result.Reason
		}

		final.HitWords = append(final.HitWords, r.Result.HitWords...)
		if r.Result.SafetyType != "" {
			final.SafetyType = r.Result.SafetyType
		}
		if r.Result.AttackType != "" {
			final.AttackType = r.Result.AttackType
		}
	}

	return final
}
