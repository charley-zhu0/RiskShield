package repository

import (
	"context"

	"github.com/charley/riskshield/internal/domain"
)

// InterventionRepository 干预库仓储接口
type InterventionRepository interface {
	FindAll(ctx context.Context, filters *domain.InterventionQueryDTO) ([]*domain.Intervention, int64, error)
	FindByID(ctx context.Context, id uint) (*domain.Intervention, error)
	Create(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error)
	Update(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error)
	Delete(ctx context.Context, id uint) error
}
