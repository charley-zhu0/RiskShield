package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/charley/riskshield/internal/domain"
	"gorm.io/gorm"
)

type sensitiveWordRepositoryImpl struct {
	db *gorm.DB
}

// NewSensitiveWordRepository 创建敏感词仓储实例
func NewSensitiveWordRepository(db *gorm.DB) SensitiveWordRepository {
	return &sensitiveWordRepositoryImpl{db: db}
}

func (r *sensitiveWordRepositoryImpl) FindAll(ctx context.Context, filters *domain.SensitiveWordQueryDTO) ([]*domain.SensitiveWord, int64, error) {
	var words []*domain.SensitiveWord
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.SensitiveWord{})

	// 应用过滤条件
	if filters.Word != "" {
		// 模糊查询 - 转义特殊字符
		escapedWord := EscapeLikeString(filters.Word)
		query = query.Where("word LIKE ?", "%"+escapedWord+"%")
	}

	if filters.FirstLabel != "" {
		query = query.Where("firstLabel = ?", filters.FirstLabel)
	}

	if filters.SecondLabel != "" {
		query = query.Where("secondLabel = ?", filters.SecondLabel)
	}

	if filters.QueryDeal != nil && len(*filters.QueryDeal) > 0 {
		query = query.Where("queryDeal IN ?", *filters.QueryDeal)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count sensitive words: %w", err)
	}

	// 分页查询
	offset := (filters.Page - 1) * filters.Size
	if err := query.Offset(offset).Limit(filters.Size).Find(&words).Error; err != nil {
		return nil, 0, fmt.Errorf("find sensitive words: %w", err)
	}

	return words, total, nil
}

func (r *sensitiveWordRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.SensitiveWord, error) {
	var word domain.SensitiveWord
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&word).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find sensitive word by id: %w", err)
	}
	return &word, nil
}

func (r *sensitiveWordRepositoryImpl) CreateBatch(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error) {
	words := make([]*domain.SensitiveWord, 0, len(data.Words))

	// 使用事务批量创建
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, w := range data.Words {
			word := &domain.SensitiveWord{
				Word:        w,
				FirstLabel:  data.FirstLabel,
				SecondLabel: data.SecondLabel,
				ThirdLabel:  data.ThirdLabel,
				QueryDeal:   data.QueryDeal,
				MatchedID:   data.MatchedID,
				TenantID:    "1", // 默认租户
				App:         "default",
				Location:    "default",
				Source:      2,
				CreateBy:    "anonym",
				UpdateBy:    "anonym",
			}

			if err := tx.Create(word).Error; err != nil {
				return fmt.Errorf("create sensitive word: %w", err)
			}

			words = append(words, word)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return words, nil
}

func (r *sensitiveWordRepositoryImpl) Update(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error) {
	// 检查记录是否存在
	word, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if word == nil {
		return nil, fmt.Errorf("sensitive word not found")
	}

	// 更新字段
	word.Word = data.Word
	word.FirstLabel = data.FirstLabel
	word.SecondLabel = data.SecondLabel
	word.ThirdLabel = data.ThirdLabel
	word.QueryDeal = data.QueryDeal
	word.MatchedID = data.MatchedID
	word.UpdateBy = "anonym"

	if err := r.db.WithContext(ctx).Save(word).Error; err != nil {
		return nil, fmt.Errorf("update sensitive word: %w", err)
	}

	return word, nil
}

func (r *sensitiveWordRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.SensitiveWord{})
	if result.Error != nil {
		return fmt.Errorf("delete sensitive word: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("sensitive word not found")
	}

	return nil
}
