package service

import (
	"context"
	"errors"
	"testing"

	"github.com/charley/riskshield/internal/domain"
)

// MockLabelRepository 模拟 LabelRepository
type MockLabelRepository struct {
	FindAllFunc  func(ctx context.Context, filters *domain.LabelQueryDTO) ([]*domain.Label, int64, error)
	FindByIDFunc func(ctx context.Context, id uint) (*domain.Label, error)
	CreateFunc   func(ctx context.Context, data *domain.LabelAddDTO) (*domain.Label, error)
	UpdateFunc   func(ctx context.Context, id uint, data *domain.LabelEditDTO) (*domain.Label, error)
	DeleteFunc   func(ctx context.Context, id uint) error
}

func (m *MockLabelRepository) FindAll(ctx context.Context, filters *domain.LabelQueryDTO) ([]*domain.Label, int64, error) {
	return m.FindAllFunc(ctx, filters)
}

func (m *MockLabelRepository) FindByID(ctx context.Context, id uint) (*domain.Label, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *MockLabelRepository) FindBySID(ctx context.Context, sid string) (*domain.Label, error) {
	return nil, nil
}

func (m *MockLabelRepository) Create(ctx context.Context, data *domain.LabelAddDTO) (*domain.Label, error) {
	return m.CreateFunc(ctx, data)
}

func (m *MockLabelRepository) Update(ctx context.Context, id uint, data *domain.LabelEditDTO) (*domain.Label, error) {
	return m.UpdateFunc(ctx, id, data)
}

func (m *MockLabelRepository) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func TestLabelService_Query(t *testing.T) {
	tests := []struct {
		name      string
		dto       *domain.LabelQueryDTO
		mockSetup func(*MockLabelRepository)
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "成功查询",
			dto: &domain.LabelQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock *MockLabelRepository) {
				mock.FindAllFunc = func(ctx context.Context, filters *domain.LabelQueryDTO) ([]*domain.Label, int64, error) {
					return []*domain.Label{
						{ID: 1, Title: "标签1"},
						{ID: 2, Title: "标签2"},
					}, 2, nil
				}
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "Repository 错误",
			dto: &domain.LabelQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock *MockLabelRepository) {
				mock.FindAllFunc = func(ctx context.Context, filters *domain.LabelQueryDTO) ([]*domain.Label, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockLabelRepository{}
			tt.mockSetup(mock)

			service := NewLabelService(mock)
			resp, err := service.Query(context.Background(), tt.dto)

			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(resp.Data) != tt.wantCount {
					t.Errorf("Query() got %d labels, want %d", len(resp.Data), tt.wantCount)
				}
				if resp.Total != tt.wantTotal {
					t.Errorf("Query() got total %d, want %d", resp.Total, tt.wantTotal)
				}
			}
		})
	}
}

func TestLabelService_Add(t *testing.T) {
	tests := []struct {
		name      string
		dto       *domain.LabelAddDTO
		mockSetup func(*MockLabelRepository)
		wantErr   bool
	}{
		{
			name: "成功添加",
			dto: &domain.LabelAddDTO{
				PID:   "zhinaolabel",
				Title: "新标签",
			},
			mockSetup: func(mock *MockLabelRepository) {
				mock.CreateFunc = func(ctx context.Context, data *domain.LabelAddDTO) (*domain.Label, error) {
					return &domain.Label{
						ID:    1,
						Title: data.Title,
						PID:   data.PID,
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "Repository 错误",
			dto: &domain.LabelAddDTO{
				PID:   "zhinaolabel",
				Title: "新标签",
			},
			mockSetup: func(mock *MockLabelRepository) {
				mock.CreateFunc = func(ctx context.Context, data *domain.LabelAddDTO) (*domain.Label, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockLabelRepository{}
			tt.mockSetup(mock)

			service := NewLabelService(mock)
			label, err := service.Add(context.Background(), tt.dto)

			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && label == nil {
				t.Error("Add() expected non-nil label")
			}
		})
	}
}

func TestLabelService_Edit(t *testing.T) {
	tests := []struct {
		name      string
		dto       *domain.LabelEditDTO
		mockSetup func(*MockLabelRepository)
		wantErr   bool
	}{
		{
			name: "成功编辑",
			dto: &domain.LabelEditDTO{
				ID:    1,
				PID:   "zhinaolabel",
				Title: "更新后的标签",
			},
			mockSetup: func(mock *MockLabelRepository) {
				mock.UpdateFunc = func(ctx context.Context, id uint, data *domain.LabelEditDTO) (*domain.Label, error) {
					return &domain.Label{
						ID:    id,
						Title: data.Title,
						PID:   data.PID,
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "Repository 错误",
			dto: &domain.LabelEditDTO{
				ID:    1,
				PID:   "zhinaolabel",
				Title: "更新后的标签",
			},
			mockSetup: func(mock *MockLabelRepository) {
				mock.UpdateFunc = func(ctx context.Context, id uint, data *domain.LabelEditDTO) (*domain.Label, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockLabelRepository{}
			tt.mockSetup(mock)

			service := NewLabelService(mock)
			label, err := service.Edit(context.Background(), tt.dto)

			if (err != nil) != tt.wantErr {
				t.Errorf("Edit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && label == nil {
				t.Error("Edit() expected non-nil label")
			}
		})
	}
}

func TestLabelService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		dto       *domain.LabelDeleteDTO
		mockSetup func(*MockLabelRepository)
		wantErr   bool
	}{
		{
			name: "成功删除",
			dto: &domain.LabelDeleteDTO{
				ID: 1,
			},
			mockSetup: func(mock *MockLabelRepository) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "Repository 错误",
			dto: &domain.LabelDeleteDTO{
				ID: 1,
			},
			mockSetup: func(mock *MockLabelRepository) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockLabelRepository{}
			tt.mockSetup(mock)

			service := NewLabelService(mock)
			err := service.Delete(context.Background(), tt.dto)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
