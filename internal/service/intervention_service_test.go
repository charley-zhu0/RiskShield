package service

import (
	"context"
	"errors"
	"testing"

	"github.com/charley/riskshield/internal/domain"
)

// MockInterventionRepository 用于测试的 Mock Repository
type MockInterventionRepository struct {
	FindAllFunc  func(ctx context.Context, filters *domain.InterventionQueryDTO) ([]*domain.Intervention, int64, error)
	FindByIDFunc func(ctx context.Context, id uint) (*domain.Intervention, error)
	CreateFunc   func(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error)
	UpdateFunc   func(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error)
	DeleteFunc   func(ctx context.Context, id uint) error
}

func (m *MockInterventionRepository) FindAll(ctx context.Context, filters *domain.InterventionQueryDTO) ([]*domain.Intervention, int64, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, filters)
	}
	return nil, 0, nil
}

func (m *MockInterventionRepository) FindByID(ctx context.Context, id uint) (*domain.Intervention, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockInterventionRepository) Create(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, data)
	}
	return nil, nil
}

func (m *MockInterventionRepository) Update(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, data)
	}
	return nil, nil
}

func (m *MockInterventionRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestInterventionService_Query(t *testing.T) {
	tests := []struct {
		name      string
		filters   *domain.InterventionQueryDTO
		mockSetup func(*MockInterventionRepository)
		wantErr   bool
	}{
		{
			name: "成功查询",
			filters: &domain.InterventionQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock *MockInterventionRepository) {
				mock.FindAllFunc = func(ctx context.Context, filters *domain.InterventionQueryDTO) ([]*domain.Intervention, int64, error) {
					return []*domain.Intervention{
						{ID: 1, Query: "test1", Answer: "answer1"},
						{ID: 2, Query: "test2", Answer: "answer2"},
					}, 2, nil
				}
			},
			wantErr: false,
		},
		{
			name: "Repository 错误",
			filters: &domain.InterventionQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock *MockInterventionRepository) {
				mock.FindAllFunc = func(ctx context.Context, filters *domain.InterventionQueryDTO) ([]*domain.Intervention, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockInterventionRepository{}
			tt.mockSetup(mock)

			service := NewInterventionService(mock)
			response, err := service.Query(context.Background(), tt.filters)

			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && response == nil {
				t.Error("Query() expected non-nil response")
			}
		})
	}
}

func TestInterventionService_Add(t *testing.T) {
	tests := []struct {
		name      string
		data      *domain.InterventionAddDTO
		mockSetup func(*MockInterventionRepository)
		wantErr   bool
	}{
		{
			name: "成功添加",
			data: &domain.InterventionAddDTO{
				Query:  "test query",
				Answer: "test answer",
				Source: 1,
			},
			mockSetup: func(mock *MockInterventionRepository) {
				mock.CreateFunc = func(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error) {
					return &domain.Intervention{
						ID:     1,
						Query:  data.Query,
						Answer: data.Answer,
						Source: data.Source,
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "Repository 错误",
			data: &domain.InterventionAddDTO{
				Query:  "test query",
				Answer: "test answer",
				Source: 1,
			},
			mockSetup: func(mock *MockInterventionRepository) {
				mock.CreateFunc = func(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockInterventionRepository{}
			tt.mockSetup(mock)

			service := NewInterventionService(mock)
			intervention, err := service.Add(context.Background(), tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && intervention == nil {
				t.Error("Add() expected non-nil intervention")
			}
		})
	}
}

func TestInterventionService_Edit(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		data      *domain.InterventionEditDTO
		mockSetup func(*MockInterventionRepository)
		wantErr   bool
	}{
		{
			name: "成功编辑",
			id:   1,
			data: &domain.InterventionEditDTO{
				ID:     1,
				Query:  "updated query",
				Answer: "updated answer",
				Source: 2,
			},
			mockSetup: func(mock *MockInterventionRepository) {
				mock.UpdateFunc = func(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error) {
					return &domain.Intervention{
						ID:     id,
						Query:  data.Query,
						Answer: data.Answer,
						Source: data.Source,
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "记录不存在",
			id:   999,
			data: &domain.InterventionEditDTO{
				ID:     999,
				Query:  "updated query",
				Answer: "updated answer",
				Source: 2,
			},
			mockSetup: func(mock *MockInterventionRepository) {
				mock.UpdateFunc = func(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error) {
					return nil, errors.New("intervention not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockInterventionRepository{}
			tt.mockSetup(mock)

			service := NewInterventionService(mock)
			intervention, err := service.Edit(context.Background(), tt.id, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Edit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && intervention == nil {
				t.Error("Edit() expected non-nil intervention")
			}
		})
	}
}

func TestInterventionService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func(*MockInterventionRepository)
		wantErr   bool
	}{
		{
			name: "成功删除",
			id:   1,
			mockSetup: func(mock *MockInterventionRepository) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "记录不存在",
			id:   999,
			mockSetup: func(mock *MockInterventionRepository) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return errors.New("intervention not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockInterventionRepository{}
			tt.mockSetup(mock)

			service := NewInterventionService(mock)
			err := service.Delete(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
