package domain

import "context"

// DetectionResult 检测结果
type DetectionResult struct {
	ShouldReject bool
	Reason       string
	HitWords     []string
	SafetyType   string
	AttackType   string
}

// Detector 检测器接口
type Detector interface {
	Detect(ctx context.Context, content string) (*DetectionResult, error)
	Name() string
}
