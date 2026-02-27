package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/charley/riskshield/internal/domain"
)

// MockRiskEngineRepository 模拟 Repository
type MockRiskEngineRepository struct {
	FindAllFunc  func(ctx context.Context, filters *domain.RiskEngineQueryDTO) ([]*domain.RiskEngine, int64, error)
	FindByIDFunc func(ctx context.Context, id uint) (*domain.RiskEngine, error)
	CreateFunc   func(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error)
	UpdateFunc   func(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error)
	DeleteFunc   func(ctx context.Context, id uint) error
}

func (m *MockRiskEngineRepository) FindAll(ctx context.Context, filters *domain.RiskEngineQueryDTO) ([]*domain.RiskEngine, int64, error) {
	return m.FindAllFunc(ctx, filters)
}

func (m *MockRiskEngineRepository) FindByID(ctx context.Context, id uint) (*domain.RiskEngine, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *MockRiskEngineRepository) Create(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error) {
	return m.CreateFunc(ctx, data)
}

func (m *MockRiskEngineRepository) Update(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error) {
	return m.UpdateFunc(ctx, id, data)
}

func (m *MockRiskEngineRepository) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func TestRiskEngineService_Query(t *testing.T) {
	tests := []struct {
		name      string
		filters   *domain.RiskEngineQueryDTO
		mockSetup func(*MockRiskEngineRepository)
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "成功查询",
			filters: &domain.RiskEngineQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock *MockRiskEngineRepository) {
				mock.FindAllFunc = func(ctx context.Context, filters *domain.RiskEngineQueryDTO) ([]*domain.RiskEngine, int64, error) {
					return []*domain.RiskEngine{
						{ID: 1, App: "app1", Location: "loc1", FirstLabel: "first1", SecondLabel: "second1", ThirdLabel: "third1", QueryDeal: 0, Source: 1},
						{ID: 2, App: "app2", Location: "loc2", FirstLabel: "first2", SecondLabel: "second2", ThirdLabel: "third2", QueryDeal: 1, Source: 2},
					}, 2, nil
				}
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "Repository 错误",
			filters: &domain.RiskEngineQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock *MockRiskEngineRepository) {
				mock.FindAllFunc = func(ctx context.Context, filters *domain.RiskEngineQueryDTO) ([]*domain.RiskEngine, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRiskEngineRepository{}
			tt.mockSetup(mockRepo)

			service := NewRiskEngineService(mockRepo)
			response, err := service.Query(context.Background(), tt.filters)

			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(response.Data) != tt.wantCount {
					t.Errorf("Query() got %d engines, want %d", len(response.Data), tt.wantCount)
				}
				if response.Total != tt.wantTotal {
					t.Errorf("Query() got total %d, want %d", response.Total, tt.wantTotal)
				}
			}
		})
	}
}

func TestRiskEngineService_Add(t *testing.T) {
	tests := []struct {
		name      string
		data      *domain.RiskEngineAddDTO
		mockSetup func(*MockRiskEngineRepository)
		wantErr   bool
	}{
		{
			name: "成功添加",
			data: &domain.RiskEngineAddDTO{
				App:         "test_app",
				Location:    "test_location",
				FirstLabel:  "first",
				SecondLabel: "second",
				ThirdLabel:  "third",
				QueryDeal:   0,
				Source:      1,
			},
			mockSetup: func(mock *MockRiskEngineRepository) {
				mock.CreateFunc = func(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error) {
					return &domain.RiskEngine{
						ID:          1,
						App:         data.App,
						Location:    data.Location,
						FirstLabel:  data.FirstLabel,
						SecondLabel: data.SecondLabel,
						ThirdLabel:  data.ThirdLabel,
						QueryDeal:   data.QueryDeal,
						Source:      data.Source,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "Repository 错误",
			data: &domain.RiskEngineAddDTO{
				App:         "test_app",
				Location:    "test_location",
				FirstLabel:  "first",
				SecondLabel: "second",
				ThirdLabel:  "third",
				QueryDeal:   0,
				Source:      1,
			},
			mockSetup: func(mock *MockRiskEngineRepository) {
				mock.CreateFunc = func(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRiskEngineRepository{}
			tt.mockSetup(mockRepo)

			service := NewRiskEngineService(mockRepo)
			engine, err := service.Add(context.Background(), tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && engine == nil {
				t.Error("Add() expected non-nil engine")
			}
		})
	}
}

func TestRiskEngineService_Edit(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		data      *domain.RiskEngineEditDTO
		mockSetup func(*MockRiskEngineRepository)
		wantErr   bool
	}{
		{
			name: "成功编辑",
			id:   1,
			data: &domain.RiskEngineEditDTO{
				ID:          1,
				App:         "updated_app",
				Location:    "updated_location",
				FirstLabel:  "updated_first",
				SecondLabel: "updated_second",
				ThirdLabel:  "updated_third",
				QueryDeal:   1,
				Source:      2,
			},
			mockSetup: func(mock *MockRiskEngineRepository) {
				mock.UpdateFunc = func(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error) {
					return &domain.RiskEngine{
						ID:          id,
						App:         data.App,
						Location:    data.Location,
						FirstLabel:  data.FirstLabel,
						SecondLabel: data.SecondLabel,
						ThirdLabel:  data.ThirdLabel,
						QueryDeal:   data.QueryDeal,
						Source:      data.Source,
						UpdatedAt:   time.Now(),
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "Repository 错误",
			id:   1,
			data: &domain.RiskEngineEditDTO{
				ID:          1,
				App:         "updated_app",
				Location:    "updated_location",
				FirstLabel:  "updated_first",
				SecondLabel: "updated_second",
				ThirdLabel:  "updated_third",
				QueryDeal:   1,
				Source:      2,
			},
			mockSetup: func(mock *MockRiskEngineRepository) {
				mock.UpdateFunc = func(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRiskEngineRepository{}
			tt.mockSetup(mockRepo)

			service := NewRiskEngineService(mockRepo)
			engine, err := service.Edit(context.Background(), tt.id, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Edit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && engine == nil {
				t.Error("Edit() expected non-nil engine")
			}
		})
	}
}

func TestRiskEngineService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func(*MockRiskEngineRepository)
		wantErr   bool
	}{
		{
			name: "成功删除",
			id:   1,
			mockSetup: func(mock *MockRiskEngineRepository) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "Repository 错误",
			id:   1,
			mockSetup: func(mock *MockRiskEngineRepository) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRiskEngineRepository{}
			tt.mockSetup(mockRepo)

			service := NewRiskEngineService(mockRepo)
			err := service.Delete(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
