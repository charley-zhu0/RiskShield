package repository

import (
	"context"

	"github.com/charley/riskshield/internal/domain"
)

// LabelRepository 标签仓储接口
type LabelRepository interface {
	FindAll(ctx context.Context, filters *domain.LabelQueryDTO) ([]*domain.Label, int64, error)
	FindByID(ctx context.Context, id uint) (*domain.Label, error)
	FindBySID(ctx context.Context, sid string) (*domain.Label, error)
	Create(ctx context.Context, data *domain.LabelAddDTO) (*domain.Label, error)
	Update(ctx context.Context, id uint, data *domain.LabelEditDTO) (*domain.Label, error)
	Delete(ctx context.Context, id uint) error
}
