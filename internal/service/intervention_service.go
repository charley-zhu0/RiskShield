package service

import (
	"context"

	"github.com/charley/riskshield/internal/domain"
)

// InterventionService 干预库服务接口
type InterventionService interface {
	Query(ctx context.Context, filters *domain.InterventionQueryDTO) (*domain.InterventionQueryResponse, error)
	Add(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error)
	Edit(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error)
	Delete(ctx context.Context, id uint) error
}
