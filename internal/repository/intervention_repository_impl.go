package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/charley/riskshield/internal/domain"
	"gorm.io/gorm"
)

type interventionRepositoryImpl struct {
	db *gorm.DB
}

// NewInterventionRepository 创建干预库仓储实例
func NewInterventionRepository(db *gorm.DB) InterventionRepository {
	return &interventionRepositoryImpl{db: db}
}

func (r *interventionRepositoryImpl) FindAll(ctx context.Context, filters *domain.InterventionQueryDTO) ([]*domain.Intervention, int64, error) {
	var interventions []*domain.Intervention
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Intervention{})

	// 应用过滤条件
	if filters.Query != "" {
		if filters.Fuzzy {
			// 模糊查询 - 转义特殊字符
			escapedQuery := EscapeLikeString(filters.Query)
			query = query.Where("query LIKE ?", "%"+escapedQuery+"%")
		} else {
			// 精确查询（使用 queryHash）
			queryHash := domain.CalculateQueryHash(filters.Query)
			query = query.Where("queryHash = ?", queryHash)
		}
	}

	if filters.Source != nil {
		query = query.Where("source = ?", *filters.Source)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count interventions: %w", err)
	}

	// 分页查询
	offset := (filters.Page - 1) * filters.Size
	if err := query.Offset(offset).Limit(filters.Size).Find(&interventions).Error; err != nil {
		return nil, 0, fmt.Errorf("find interventions: %w", err)
	}

	return interventions, total, nil
}

func (r *interventionRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Intervention, error) {
	var intervention domain.Intervention
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&intervention).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find intervention by id: %w", err)
	}
	return &intervention, nil
}

func (r *interventionRepositoryImpl) Create(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error) {
	intervention := &domain.Intervention{
		Query:     data.Query,
		Answer:    data.Answer,
		QueryHash: domain.CalculateQueryHash(data.Query),
		Source:    data.Source,
	}

	if err := r.db.WithContext(ctx).Create(intervention).Error; err != nil {
		return nil, fmt.Errorf("create intervention: %w", err)
	}

	return intervention, nil
}

func (r *interventionRepositoryImpl) Update(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error) {
	// 检查记录是否存在
	intervention, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if intervention == nil {
		return nil, fmt.Errorf("intervention not found")
	}

	// 更新字段
	intervention.Query = data.Query
	intervention.Answer = data.Answer
	intervention.QueryHash = domain.CalculateQueryHash(data.Query)
	intervention.Source = data.Source

	if err := r.db.WithContext(ctx).Save(intervention).Error; err != nil {
		return nil, fmt.Errorf("update intervention: %w", err)
	}

	return intervention, nil
}

func (r *interventionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Intervention{})
	if result.Error != nil {
		return fmt.Errorf("delete intervention: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("intervention not found")
	}

	return nil
}
