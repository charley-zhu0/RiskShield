package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/charley/riskshield/internal/domain"
	"gorm.io/gorm"
)

type riskEngineRepositoryImpl struct {
	db *gorm.DB
}

// NewRiskEngineRepository 创建防护策略仓储实例
func NewRiskEngineRepository(db *gorm.DB) RiskEngineRepository {
	return &riskEngineRepositoryImpl{db: db}
}

func (r *riskEngineRepositoryImpl) FindAll(ctx context.Context, filters *domain.RiskEngineQueryDTO) ([]*domain.RiskEngine, int64, error) {
	var engines []*domain.RiskEngine
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.RiskEngine{})

	// 应用过滤条件
	if filters.ThirdLabel != nil && *filters.ThirdLabel != "" {
		query = query.Where("thirdLabel = ?", *filters.ThirdLabel)
	}

	if filters.Locale != nil && *filters.Locale != "" {
		query = query.Where("location = ?", *filters.Locale)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count risk engines: %w", err)
	}

	// 分页查询
	offset := (filters.Page - 1) * filters.Size
	if err := query.Offset(offset).Limit(filters.Size).Find(&engines).Error; err != nil {
		return nil, 0, fmt.Errorf("find risk engines: %w", err)
	}

	return engines, total, nil
}

func (r *riskEngineRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.RiskEngine, error) {
	var engine domain.RiskEngine
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&engine).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find risk engine by id: %w", err)
	}
	return &engine, nil
}

func (r *riskEngineRepositoryImpl) Create(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error) {
	engine := &domain.RiskEngine{
		App:         data.App,
		Location:    data.Location,
		FirstLabel:  data.FirstLabel,
		SecondLabel: data.SecondLabel,
		ThirdLabel:  data.ThirdLabel,
		QueryDeal:   data.QueryDeal,
		Source:      data.Source,
	}

	if err := r.db.WithContext(ctx).Create(engine).Error; err != nil {
		return nil, fmt.Errorf("create risk engine: %w", err)
	}

	return engine, nil
}

func (r *riskEngineRepositoryImpl) Update(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error) {
	// 检查记录是否存在
	engine, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if engine == nil {
		return nil, fmt.Errorf("risk engine not found")
	}

	// 更新字段
	engine.App = data.App
	engine.Location = data.Location
	engine.FirstLabel = data.FirstLabel
	engine.SecondLabel = data.SecondLabel
	engine.ThirdLabel = data.ThirdLabel
	engine.QueryDeal = data.QueryDeal
	engine.Source = data.Source

	if err := r.db.WithContext(ctx).Save(engine).Error; err != nil {
		return nil, fmt.Errorf("update risk engine: %w", err)
	}

	return engine, nil
}

func (r *riskEngineRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.RiskEngine{})
	if result.Error != nil {
		return fmt.Errorf("delete risk engine: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("risk engine not found")
	}

	return nil
}
