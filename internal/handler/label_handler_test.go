package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charley/riskshield/internal/domain"
	"github.com/gin-gonic/gin"
)

// MockLabelService 模拟 LabelService
type MockLabelService struct {
	QueryFunc  func(ctx context.Context, dto *domain.LabelQueryDTO) (*domain.LabelQueryResponse, error)
	AddFunc    func(ctx context.Context, dto *domain.LabelAddDTO) (*domain.Label, error)
	EditFunc   func(ctx context.Context, dto *domain.LabelEditDTO) (*domain.Label, error)
	DeleteFunc func(ctx context.Context, dto *domain.LabelDeleteDTO) error
}

func (m *MockLabelService) Query(ctx context.Context, dto *domain.LabelQueryDTO) (*domain.LabelQueryResponse, error) {
	return m.QueryFunc(ctx, dto)
}

func (m *MockLabelService) Add(ctx context.Context, dto *domain.LabelAddDTO) (*domain.Label, error) {
	return m.AddFunc(ctx, dto)
}

func (m *MockLabelService) Edit(ctx context.Context, dto *domain.LabelEditDTO) (*domain.Label, error) {
	return m.EditFunc(ctx, dto)
}

func (m *MockLabelService) Delete(ctx context.Context, dto *domain.LabelDeleteDTO) error {
	return m.DeleteFunc(ctx, dto)
}

func TestLabelHandler_Query(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		query      string
		mockSetup  func(*MockLabelService)
		wantStatus int
	}{
		{
			name:  "成功查询",
			query: "?page=1&size=10",
			mockSetup: func(mock *MockLabelService) {
				mock.QueryFunc = func(ctx context.Context, dto *domain.LabelQueryDTO) (*domain.LabelQueryResponse, error) {
					return &domain.LabelQueryResponse{
						Data:  []domain.Label{{ID: 1, Title: "标签1"}},
						Total: 1,
						Page:  1,
						Size:  10,
					}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "参数验证失败",
			query: "?page=0&size=10",
			mockSetup: func(mock *MockLabelService) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "Service 错误",
			query: "?page=1&size=10",
			mockSetup: func(mock *MockLabelService) {
				mock.QueryFunc = func(ctx context.Context, dto *domain.LabelQueryDTO) (*domain.LabelQueryResponse, error) {
					return nil, errors.New("service error")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockLabelService{}
			tt.mockSetup(mock)

			handler := NewLabelHandler(mock)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/label/query"+tt.query, nil)

			handler.Query(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Query() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestLabelHandler_Add(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		mockSetup  func(*MockLabelService)
		wantStatus int
	}{
		{
			name: "成功添加",
			body: domain.LabelAddDTO{
				PID:   "zhinaolabel",
				Title: "新标签",
			},
			mockSetup: func(mock *MockLabelService) {
				mock.AddFunc = func(ctx context.Context, dto *domain.LabelAddDTO) (*domain.Label, error) {
					return &domain.Label{ID: 1, Title: dto.Title}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "参数验证失败",
			body: domain.LabelAddDTO{
				PID: "zhinaolabel",
				// Title 缺失
			},
			mockSetup:  func(mock *MockLabelService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Service 错误",
			body: domain.LabelAddDTO{
				PID:   "zhinaolabel",
				Title: "新标签",
			},
			mockSetup: func(mock *MockLabelService) {
				mock.AddFunc = func(ctx context.Context, dto *domain.LabelAddDTO) (*domain.Label, error) {
					return nil, errors.New("service error")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockLabelService{}
			tt.mockSetup(mock)

			handler := NewLabelHandler(mock)

			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/label/add", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Add(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Add() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestLabelHandler_Edit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		mockSetup  func(*MockLabelService)
		wantStatus int
	}{
		{
			name: "成功编辑",
			body: domain.LabelEditDTO{
				ID:    1,
				PID:   "zhinaolabel",
				Title: "更新后的标签",
			},
			mockSetup: func(mock *MockLabelService) {
				mock.EditFunc = func(ctx context.Context, dto *domain.LabelEditDTO) (*domain.Label, error) {
					return &domain.Label{ID: dto.ID, Title: dto.Title}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "参数验证失败",
			body: domain.LabelEditDTO{
				// ID 缺失
				PID:   "zhinaolabel",
				Title: "更新后的标签",
			},
			mockSetup:  func(mock *MockLabelService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockLabelService{}
			tt.mockSetup(mock)

			handler := NewLabelHandler(mock)

			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/label/edit", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Edit(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Edit() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestLabelHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		mockSetup  func(*MockLabelService)
		wantStatus int
	}{
		{
			name: "成功删除",
			body: domain.LabelDeleteDTO{
				ID: 1,
			},
			mockSetup: func(mock *MockLabelService) {
				mock.DeleteFunc = func(ctx context.Context, dto *domain.LabelDeleteDTO) error {
					return nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "参数验证失败",
			body: domain.LabelDeleteDTO{
				// ID 缺失
			},
			mockSetup:  func(mock *MockLabelService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockLabelService{}
			tt.mockSetup(mock)

			handler := NewLabelHandler(mock)

			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/label/delete", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Delete(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Delete() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
