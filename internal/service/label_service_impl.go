package service

import (
	"context"
	"fmt"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/repository"
)

type labelService struct {
	labelRepo repository.LabelRepository
}

// NewLabelService 创建标签服务实例
func NewLabelService(labelRepo repository.LabelRepository) LabelService {
	return &labelService{
		labelRepo: labelRepo,
	}
}

// Query 查询标签列表
func (s *labelService) Query(ctx context.Context, dto *domain.LabelQueryDTO) (*domain.LabelQueryResponse, error) {
	labels, total, err := s.labelRepo.FindAll(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to query labels: %w", err)
	}

	// 转换为响应格式
	data := make([]domain.Label, len(labels))
	for i, label := range labels {
		data[i] = *label
	}

	return &domain.LabelQueryResponse{
		Data:  data,
		Total: total,
		Page:  dto.Page,
		Size:  dto.Size,
	}, nil
}

// Add 添加标签
func (s *labelService) Add(ctx context.Context, dto *domain.LabelAddDTO) (*domain.Label, error) {
	label, err := s.labelRepo.Create(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to add label: %w", err)
	}
	return label, nil
}

// Edit 编辑标签
func (s *labelService) Edit(ctx context.Context, dto *domain.LabelEditDTO) (*domain.Label, error) {
	label, err := s.labelRepo.Update(ctx, dto.ID, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to edit label: %w", err)
	}
	return label, nil
}

// Delete 删除标签
func (s *labelService) Delete(ctx context.Context, dto *domain.LabelDeleteDTO) error {
	if err := s.labelRepo.Delete(ctx, dto.ID); err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}
	return nil
}
