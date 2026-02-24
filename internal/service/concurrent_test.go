package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/charley/riskshield/internal/domain"
)

type mockDetector struct {
	name   string
	result *domain.DetectionResult
	err    error
	delay  time.Duration
}

func (m *mockDetector) Detect(ctx context.Context, content string) (*domain.DetectionResult, error) {
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return m.result, m.err
}

func (m *mockDetector) Name() string {
	return m.name
}

func TestConcurrentExecutor(t *testing.T) {
	t.Run("所有检测器成功", func(t *testing.T) {
		detectors := []domain.Detector{
			&mockDetector{
				name:   "detector1",
				result: &domain.DetectionResult{ShouldReject: false},
			},
			&mockDetector{
				name:   "detector2",
				result: &domain.DetectionResult{ShouldReject: false},
			},
		}

		executor := NewConcurrentExecutor(detectors, 300*time.Millisecond)
		results := executor.Execute(context.Background(), "test")

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("超时忽略慢检测器", func(t *testing.T) {
		detectors := []domain.Detector{
			&mockDetector{
				name:   "fast",
				result: &domain.DetectionResult{ShouldReject: false},
				delay:  10 * time.Millisecond,
			},
			&mockDetector{
				name:   "slow",
				result: &domain.DetectionResult{ShouldReject: false},
				delay:  500 * time.Millisecond,
			},
		}

		executor := NewConcurrentExecutor(detectors, 100*time.Millisecond)
		results := executor.Execute(context.Background(), "test")

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if results[0].DetectorName != "fast" {
			t.Errorf("expected fast detector, got %s", results[0].DetectorName)
		}
	})

	t.Run("检测器返回错误", func(t *testing.T) {
		detectors := []domain.Detector{
			&mockDetector{
				name: "error",
				err:  errors.New("detector error"),
			},
		}

		executor := NewConcurrentExecutor(detectors, 300*time.Millisecond)
		results := executor.Execute(context.Background(), "test")

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if results[0].Error == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAggregator(t *testing.T) {
	t.Run("任一拒绝则拒绝", func(t *testing.T) {
		results := []DetectionResult{
			{
				DetectorName: "detector1",
				Result:       &domain.DetectionResult{ShouldReject: false},
			},
			{
				DetectorName: "detector2",
				Result:       &domain.DetectionResult{ShouldReject: true, Reason: "违规"},
			},
		}

		aggregator := NewAggregator()
		final := aggregator.Aggregate(results)

		if final.Decision != "REJECT" {
			t.Errorf("expected REJECT, got %s", final.Decision)
		}
	})

	t.Run("全部通过则通过", func(t *testing.T) {
		results := []DetectionResult{
			{
				DetectorName: "detector1",
				Result:       &domain.DetectionResult{ShouldReject: false},
			},
			{
				DetectorName: "detector2",
				Result:       &domain.DetectionResult{ShouldReject: false},
			},
		}

		aggregator := NewAggregator()
		final := aggregator.Aggregate(results)

		if final.Decision != "PASS" {
			t.Errorf("expected PASS, got %s", final.Decision)
		}
	})

	t.Run("合并敏感词", func(t *testing.T) {
		results := []DetectionResult{
			{
				DetectorName: "sw",
				Result:       &domain.DetectionResult{ShouldReject: true, HitWords: []string{"词1", "词2"}},
			},
		}

		aggregator := NewAggregator()
		final := aggregator.Aggregate(results)

		if len(final.HitWords) != 2 {
			t.Errorf("expected 2 hit words, got %d", len(final.HitWords))
		}
	})
}
