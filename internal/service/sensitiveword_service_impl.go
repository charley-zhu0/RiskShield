package service

import (
	"context"
	"fmt"

	"github.com/charley/riskshield/internal/domain"
	"github.com/charley/riskshield/internal/repository"
)

type sensitiveWordServiceImpl struct {
	repo repository.SensitiveWordRepository
}

// NewSensitiveWordService 创建敏感词服务实例
func NewSensitiveWordService(repo repository.SensitiveWordRepository) SensitiveWordService {
	return &sensitiveWordServiceImpl{repo: repo}
}

func (s *sensitiveWordServiceImpl) Query(ctx context.Context, filters *domain.SensitiveWordQueryDTO) (*domain.SensitiveWordQueryResponse, error) {
	words, total, err := s.repo.FindAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("query sensitive words: %w", err)
	}

	// 转换为非指针切片
	list := make([]domain.SensitiveWord, len(words))
	for i, word := range words {
		list[i] = *word
	}

	return &domain.SensitiveWordQueryResponse{
		List:  list,
		Total: total,
	}, nil
}

func (s *sensitiveWordServiceImpl) Add(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error) {
	words, err := s.repo.CreateBatch(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("add sensitive words: %w", err)
	}
	return words, nil
}

func (s *sensitiveWordServiceImpl) Edit(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error) {
	word, err := s.repo.Update(ctx, id, data)
	if err != nil {
		return nil, fmt.Errorf("edit sensitive word: %w", err)
	}
	return word, nil
}

func (s *sensitiveWordServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete sensitive word: %w", err)
	}
	return nil
}
