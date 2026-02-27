package service

import (
	"context"

	"github.com/charley/riskshield/internal/domain"
)

// RiskEngineService 防护策略服务接口
type RiskEngineService interface {
	Query(ctx context.Context, filters *domain.RiskEngineQueryDTO) (*domain.RiskEngineQueryResponse, error)
	Add(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error)
	Edit(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error)
	Delete(ctx context.Context, id uint) error
}
