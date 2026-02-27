package service

import (
	"context"
	"testing"

	"github.com/charley/riskshield/internal/domain"
)

// MockSensitiveWordRepository 用于测试的 Mock Repository
type MockSensitiveWordRepository struct {
	FindAllFunc     func(ctx context.Context, filters *domain.SensitiveWordQueryDTO) ([]*domain.SensitiveWord, int64, error)
	FindByIDFunc    func(ctx context.Context, id uint) (*domain.SensitiveWord, error)
	CreateBatchFunc func(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error)
	UpdateFunc      func(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error)
	DeleteFunc      func(ctx context.Context, id uint) error
}

func (m *MockSensitiveWordRepository) FindAll(ctx context.Context, filters *domain.SensitiveWordQueryDTO) ([]*domain.SensitiveWord, int64, error) {
	return m.FindAllFunc(ctx, filters)
}

func (m *MockSensitiveWordRepository) FindByID(ctx context.Context, id uint) (*domain.SensitiveWord, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *MockSensitiveWordRepository) CreateBatch(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error) {
	return m.CreateBatchFunc(ctx, data)
}

func (m *MockSensitiveWordRepository) Update(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error) {
	return m.UpdateFunc(ctx, id, data)
}

func (m *MockSensitiveWordRepository) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func TestSensitiveWordService_Query(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		filters *domain.SensitiveWordQueryDTO
		mockFn  func(ctx context.Context, filters *domain.SensitiveWordQueryDTO) ([]*domain.SensitiveWord, int64, error)
		wantErr bool
	}{
		{
			name: "successful query",
			filters: &domain.SensitiveWordQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockFn: func(ctx context.Context, filters *domain.SensitiveWordQueryDTO) ([]*domain.SensitiveWord, int64, error) {
				return []*domain.SensitiveWord{
					{ID: 1, Word: "test"},
				}, 1, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSensitiveWordRepository{
				FindAllFunc: tt.mockFn,
			}
			service := NewSensitiveWordService(mockRepo)

			got, err := service.Query(ctx, tt.filters)

			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Error("Query() returned nil response")
			}
		})
	}
}

func TestSensitiveWordService_Add(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		data    *domain.SensitiveWordAddDTO
		mockFn  func(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error)
		wantErr bool
	}{
		{
			name: "successful add",
			data: &domain.SensitiveWordAddDTO{
				Words:       []string{"test"},
				FirstLabel:  "100",
				SecondLabel: "100001",
				QueryDeal:   1,
			},
			mockFn: func(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error) {
				return []*domain.SensitiveWord{
					{ID: 1, Word: "test"},
				}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSensitiveWordRepository{
				CreateBatchFunc: tt.mockFn,
			}
			service := NewSensitiveWordService(mockRepo)

			got, err := service.Add(ctx, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) == 0 {
				t.Error("Add() returned empty slice")
			}
		})
	}
}

func TestSensitiveWordService_Edit(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		id      uint
		data    *domain.SensitiveWordEditDTO
		mockFn  func(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error)
		wantErr bool
	}{
		{
			name: "successful edit",
			id:   1,
			data: &domain.SensitiveWordEditDTO{
				ID:          1,
				Word:        "updated",
				FirstLabel:  "100",
				SecondLabel: "100001",
				QueryDeal:   1,
			},
			mockFn: func(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error) {
				return &domain.SensitiveWord{
					ID:   id,
					Word: data.Word,
				}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSensitiveWordRepository{
				UpdateFunc: tt.mockFn,
			}
			service := NewSensitiveWordService(mockRepo)

			got, err := service.Edit(ctx, tt.id, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Edit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Error("Edit() returned nil")
			}
		})
	}
}

func TestSensitiveWordService_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		id      uint
		mockFn  func(ctx context.Context, id uint) error
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   1,
			mockFn: func(ctx context.Context, id uint) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockSensitiveWordRepository{
				DeleteFunc: tt.mockFn,
			}
			service := NewSensitiveWordService(mockRepo)

			err := service.Delete(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
