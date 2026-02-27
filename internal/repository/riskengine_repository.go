package repository

import (
	"context"

	"github.com/charley/riskshield/internal/domain"
)

// RiskEngineRepository 防护策略仓储接口
type RiskEngineRepository interface {
	FindAll(ctx context.Context, filters *domain.RiskEngineQueryDTO) ([]*domain.RiskEngine, int64, error)
	FindByID(ctx context.Context, id uint) (*domain.RiskEngine, error)
	Create(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error)
	Update(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error)
	Delete(ctx context.Context, id uint) error
}
