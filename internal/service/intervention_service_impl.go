package service

import (
	"context"
	"fmt"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/repository"
)

type interventionServiceImpl struct {
	repo repository.InterventionRepository
}

// NewInterventionService 创建干预库服务实例
func NewInterventionService(repo repository.InterventionRepository) InterventionService {
	return &interventionServiceImpl{repo: repo}
}

func (s *interventionServiceImpl) Query(ctx context.Context, filters *domain.InterventionQueryDTO) (*domain.InterventionQueryResponse, error) {
	interventions, total, err := s.repo.FindAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("query interventions: %w", err)
	}

	// 转换为非指针切片
	data := make([]domain.Intervention, len(interventions))
	for i, intervention := range interventions {
		data[i] = *intervention
	}

	return &domain.InterventionQueryResponse{
		Data:  data,
		Total: total,
		Page:  filters.Page,
		Size:  filters.Size,
	}, nil
}

func (s *interventionServiceImpl) Add(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error) {
	intervention, err := s.repo.Create(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("add intervention: %w", err)
	}
	return intervention, nil
}

func (s *interventionServiceImpl) Edit(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error) {
	intervention, err := s.repo.Update(ctx, id, data)
	if err != nil {
		return nil, fmt.Errorf("edit intervention: %w", err)
	}
	return intervention, nil
}

func (s *interventionServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete intervention: %w", err)
	}
	return nil
}
