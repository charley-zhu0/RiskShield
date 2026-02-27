package repository

import (
	"context"

	"github.com/charley/riskshield/internal/domain"
)

// SensitiveWordRepository 敏感词仓储接口
type SensitiveWordRepository interface {
	// FindAll 查询敏感词列表（支持分页和过滤）
	FindAll(ctx context.Context, filters *domain.SensitiveWordQueryDTO) ([]*domain.SensitiveWord, int64, error)

	// FindByID 根据 ID 查询敏感词
	FindByID(ctx context.Context, id uint) (*domain.SensitiveWord, error)

	// CreateBatch 批量创建敏感词
	CreateBatch(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error)

	// Update 更新敏感词
	Update(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error)

	// Delete 删除敏感词
	Delete(ctx context.Context, id uint) error
}
