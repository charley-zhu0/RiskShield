package service

import (
	"context"
	"fmt"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/repository"
)

type riskEngineServiceImpl struct {
	repo repository.RiskEngineRepository
}

// NewRiskEngineService 创建防护策略服务实例
func NewRiskEngineService(repo repository.RiskEngineRepository) RiskEngineService {
	return &riskEngineServiceImpl{repo: repo}
}

func (s *riskEngineServiceImpl) Query(ctx context.Context, filters *domain.RiskEngineQueryDTO) (*domain.RiskEngineQueryResponse, error) {
	engines, total, err := s.repo.FindAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("query risk engines: %w", err)
	}

	// 转换为非指针切片
	data := make([]domain.RiskEngine, len(engines))
	for i, engine := range engines {
		data[i] = *engine
	}

	return &domain.RiskEngineQueryResponse{
		Data:  data,
		Total: total,
		Page:  filters.Page,
		Size:  filters.Size,
	}, nil
}

func (s *riskEngineServiceImpl) Add(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error) {
	engine, err := s.repo.Create(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("add risk engine: %w", err)
	}
	return engine, nil
}

func (s *riskEngineServiceImpl) Edit(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error) {
	engine, err := s.repo.Update(ctx, id, data)
	if err != nil {
		return nil, fmt.Errorf("edit risk engine: %w", err)
	}
	return engine, nil
}

func (s *riskEngineServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete risk engine: %w", err)
	}
	return nil
}
