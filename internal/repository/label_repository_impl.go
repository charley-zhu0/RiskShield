package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/charley/riskshield/internal/domain"
	"gorm.io/gorm"
)

type labelRepository struct {
	db *gorm.DB
}

// NewLabelRepository 创建标签仓储实例
func NewLabelRepository(db *gorm.DB) LabelRepository {
	return &labelRepository{db: db}
}

// FindAll 查询标签列表
func (r *labelRepository) FindAll(ctx context.Context, filters *domain.LabelQueryDTO) ([]*domain.Label, int64, error) {
	var labels []*domain.Label
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Label{})

	// 应用过滤条件
	if filters.PID != "" {
		query = query.Where("pid = ?", filters.PID)
	}
	if filters.Title != "" {
		// 使用 LIKE 进行模糊查询，需要转义特殊字符
		escapedTitle := EscapeLikeString(filters.Title)
		query = query.Where("title LIKE ?", "%"+escapedTitle+"%")
	}
	if filters.SID != "" {
		query = query.Where("sid = ?", filters.SID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count labels: %w", err)
	}

	// 分页查询
	offset := (filters.Page - 1) * filters.Size
	if err := query.Limit(filters.Size).Offset(offset).Find(&labels).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find labels: %w", err)
	}

	return labels, total, nil
}

// FindByID 根据 ID 查询标签
func (r *labelRepository) FindByID(ctx context.Context, id uint) (*domain.Label, error) {
	var label domain.Label
	if err := r.db.WithContext(ctx).First(&label, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find label: %w", err)
	}
	return &label, nil
}

// FindBySID 根据 SID 查询标签
func (r *labelRepository) FindBySID(ctx context.Context, sid string) (*domain.Label, error) {
	var label domain.Label
	if err := r.db.WithContext(ctx).First(&label, "sid = ?", sid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find label by sid: %w", err)
	}
	return &label, nil
}

// Create 创建标签
func (r *labelRepository) Create(ctx context.Context, data *domain.LabelAddDTO) (*domain.Label, error) {
	// 生成 SID
	sid := domain.GenerateSID()

	// 查询父标签以计算 step
	var parentStep int
	if data.PID != "zhinaolabel" {
		parent, err := r.FindBySID(ctx, data.PID)
		if err != nil {
			return nil, fmt.Errorf("failed to find parent label: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("parent label not found: %s", data.PID)
		}
		parentStep = parent.Step
	}

	step := domain.CalculateStep(data.PID, parentStep)

	label := &domain.Label{
		TenantID: "1",           // 暂时硬编码
		Class:    "zhinaolabel", // 固定值
		SID:      sid,
		PID:      data.PID,
		Title:    data.Title,
		Value:    sid, // value = sid
		Step:     step,
		Source:   2, // 用户创建
	}

	if err := r.db.WithContext(ctx).Create(label).Error; err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}

	return label, nil
}

// Update 更新标签
func (r *labelRepository) Update(ctx context.Context, id uint, data *domain.LabelEditDTO) (*domain.Label, error) {
	// 检查记录是否存在
	label, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if label == nil {
		return nil, fmt.Errorf("label not found: %d", id)
	}

	// 如果 PID 改变，需要重新计算 step
	if label.PID != data.PID {
		var parentStep int
		if data.PID != "zhinaolabel" {
			parent, err := r.FindBySID(ctx, data.PID)
			if err != nil {
				return nil, fmt.Errorf("failed to find parent label: %w", err)
			}
			if parent == nil {
				return nil, fmt.Errorf("parent label not found: %s", data.PID)
			}
			parentStep = parent.Step
		}
		label.Step = domain.CalculateStep(data.PID, parentStep)
		label.PID = data.PID
	}

	// 更新字段
	label.Title = data.Title

	if err := r.db.WithContext(ctx).Save(label).Error; err != nil {
		return nil, fmt.Errorf("failed to update label: %w", err)
	}

	return label, nil
}

// Delete 删除标签
func (r *labelRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Label{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete label: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("label not found: %d", id)
	}
	return nil
}
