package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charley/riskshield/internal/domain"
	"github.com/gin-gonic/gin"
)

// MockSensitiveWordService 用于测试的 Mock Service
type MockSensitiveWordService struct {
	QueryFunc  func(ctx context.Context, filters *domain.SensitiveWordQueryDTO) (*domain.SensitiveWordQueryResponse, error)
	AddFunc    func(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error)
	EditFunc   func(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error)
	DeleteFunc func(ctx context.Context, id uint) error
}

func (m *MockSensitiveWordService) Query(ctx context.Context, filters *domain.SensitiveWordQueryDTO) (*domain.SensitiveWordQueryResponse, error) {
	return m.QueryFunc(ctx, filters)
}

func (m *MockSensitiveWordService) Add(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error) {
	return m.AddFunc(ctx, data)
}

func (m *MockSensitiveWordService) Edit(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error) {
	return m.EditFunc(ctx, id, data)
}

func (m *MockSensitiveWordService) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func TestSensitiveWordHandler_Query(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		query      string
		mockFn     func(ctx context.Context, filters *domain.SensitiveWordQueryDTO) (*domain.SensitiveWordQueryResponse, error)
		wantStatus int
	}{
		{
			name:  "successful query",
			query: "page=1&size=10",
			mockFn: func(ctx context.Context, filters *domain.SensitiveWordQueryDTO) (*domain.SensitiveWordQueryResponse, error) {
				return &domain.SensitiveWordQueryResponse{
					List:  []domain.SensitiveWord{{ID: 1, Word: "test"}},
					Total: 1,
				}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid query params",
			query:      "page=0&size=10",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSensitiveWordService{
				QueryFunc: tt.mockFn,
			}
			handler := NewSensitiveWordHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/query?"+tt.query, nil)

			handler.Query(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Query() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestSensitiveWordHandler_Add(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		mockFn     func(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error)
		wantStatus int
	}{
		{
			name: "successful add",
			body: domain.SensitiveWordAddDTO{
				Words:       []string{"test"},
				FirstLabel:  "100",
				SecondLabel: "100001",
				QueryDeal:   1,
			},
			mockFn: func(ctx context.Context, data *domain.SensitiveWordAddDTO) ([]*domain.SensitiveWord, error) {
				return []*domain.SensitiveWord{{ID: 1, Word: "test"}}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid body",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSensitiveWordService{
				AddFunc: tt.mockFn,
			}
			handler := NewSensitiveWordHandler(mockService)

			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/add", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Add(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Add() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestSensitiveWordHandler_Edit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		mockFn     func(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error)
		wantStatus int
	}{
		{
			name: "successful edit",
			body: domain.SensitiveWordEditDTO{
				ID:          1,
				Word:        "updated",
				FirstLabel:  "100",
				SecondLabel: "100001",
				QueryDeal:   1,
			},
			mockFn: func(ctx context.Context, id uint, data *domain.SensitiveWordEditDTO) (*domain.SensitiveWord, error) {
				return &domain.SensitiveWord{ID: id, Word: data.Word}, nil
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSensitiveWordService{
				EditFunc: tt.mockFn,
			}
			handler := NewSensitiveWordHandler(mockService)

			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/edit", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Edit(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Edit() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestSensitiveWordHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		mockFn     func(ctx context.Context, id uint) error
		wantStatus int
	}{
		{
			name: "successful delete",
			body: domain.SensitiveWordDeleteDTO{ID: 1},
			mockFn: func(ctx context.Context, id uint) error {
				return nil
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSensitiveWordService{
				DeleteFunc: tt.mockFn,
			}
			handler := NewSensitiveWordHandler(mockService)

			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/delete", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Delete(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Delete() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
