package service

import (
	"context"

	"github.com/charley/riskshield/internal/domain"
)

// SensitiveWordService 敏感词服务接口
type SensitiveWordService interface {
	// Query 查询敏感词列表
	Query(ctx context.Context, filters *domain.SensitiveWordQueryDTO) (*domain.SensitiveWordQueryResponse, error)

	// Add 批量添加敏感词
	Add(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error)

	// Edit 编辑敏感词
	Edit(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error)

	// Delete 删除敏感词
	Delete(ctx context.Context, id uint) error
}
