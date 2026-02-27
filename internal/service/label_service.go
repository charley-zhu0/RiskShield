package service

import (
	"context"

	"github.com/charley/riskshield/internal/domain"
)

// LabelService 标签服务接口
type LabelService interface {
	Query(ctx context.Context, dto *domain.LabelQueryDTO) (*domain.LabelQueryResponse, error)
	Add(ctx context.Context, dto *domain.LabelAddDTO) (*domain.Label, error)
	Edit(ctx context.Context, dto *domain.LabelEditDTO) (*domain.Label, error)
	Delete(ctx context.Context, dto *domain.LabelDeleteDTO) error
}
