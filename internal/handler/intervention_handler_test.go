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

// MockInterventionService 用于测试的 Mock Service
type MockInterventionService struct {
	QueryFunc  func(ctx context.Context, filters *domain.InterventionQueryDTO) (*domain.InterventionQueryResponse, error)
	AddFunc    func(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error)
	EditFunc   func(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error)
	DeleteFunc func(ctx context.Context, id uint) error
}

func (m *MockInterventionService) Query(ctx context.Context, filters *domain.InterventionQueryDTO) (*domain.InterventionQueryResponse, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, filters)
	}
	return nil, nil
}

func (m *MockInterventionService) Add(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error) {
	if m.AddFunc != nil {
		return m.AddFunc(ctx, data)
	}
	return nil, nil
}

func (m *MockInterventionService) Edit(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error) {
	if m.EditFunc != nil {
		return m.EditFunc(ctx, id, data)
	}
	return nil, nil
}

func (m *MockInterventionService) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestInterventionHandler_Query(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockInterventionService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "成功查询",
			queryParams: "?page=1&size=10",
			mockSetup: func(mock *MockInterventionService) {
				mock.QueryFunc = func(ctx context.Context, filters *domain.InterventionQueryDTO) (*domain.InterventionQueryResponse, error) {
					return &domain.InterventionQueryResponse{
						Data:  []domain.Intervention{{ID: 1, Query: "test"}},
						Total: 1,
						Page:  1,
						Size:  10,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp domain.StandardResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.ErrNo != 0 {
					t.Errorf("expected errno 0, got %d", resp.ErrNo)
				}
			},
		},
		{
			name:           "缺少必需参数",
			queryParams:    "?page=1",
			mockSetup:      func(mock *MockInterventionService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp domain.StandardResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.ErrNo == 0 {
					t.Error("expected non-zero errno")
				}
			},
		},
		{
			name:        "Service 错误",
			queryParams: "?page=1&size=10",
			mockSetup: func(mock *MockInterventionService) {
				mock.QueryFunc = func(ctx context.Context, filters *domain.InterventionQueryDTO) (*domain.InterventionQueryResponse, error) {
					return nil, errors.New("service error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp domain.StandardResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.ErrNo == 0 {
					t.Error("expected non-zero errno")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockInterventionService{}
			tt.mockSetup(mock)

			handler := NewInterventionHandler(mock)
			router := setupTestRouter()
			router.GET("/query", handler.Query)

			req := httptest.NewRequest(http.MethodGet, "/query"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestInterventionHandler_Add(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockInterventionService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "成功添加",
			requestBody: domain.InterventionAddDTO{
				Query:  "test query",
				Answer: "test answer",
				Source: 1,
			},
			mockSetup: func(mock *MockInterventionService) {
				mock.AddFunc = func(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error) {
					return &domain.Intervention{
						ID:     1,
						Query:  data.Query,
						Answer: data.Answer,
						Source: data.Source,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp domain.StandardResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.ErrNo != 0 {
					t.Errorf("expected errno 0, got %d", resp.ErrNo)
				}
			},
		},
		{
			name:           "无效的请求体",
			requestBody:    "invalid json",
			mockSetup:      func(mock *MockInterventionService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service 错误",
			requestBody: domain.InterventionAddDTO{
				Query:  "test query",
				Answer: "test answer",
				Source: 1,
			},
			mockSetup: func(mock *MockInterventionService) {
				mock.AddFunc = func(ctx context.Context, data *domain.InterventionAddDTO) (*domain.Intervention, error) {
					return nil, errors.New("service error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockInterventionService{}
			tt.mockSetup(mock)

			handler := NewInterventionHandler(mock)
			router := setupTestRouter()
			router.POST("/add", handler.Add)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestInterventionHandler_Edit(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockInterventionService)
		expectedStatus int
	}{
		{
			name: "成功编辑",
			requestBody: domain.InterventionEditDTO{
				ID:     1,
				Query:  "updated query",
				Answer: "updated answer",
				Source: 2,
			},
			mockSetup: func(mock *MockInterventionService) {
				mock.EditFunc = func(ctx context.Context, id uint, data *domain.InterventionEditDTO) (*domain.Intervention, error) {
					return &domain.Intervention{
						ID:     id,
						Query:  data.Query,
						Answer: data.Answer,
						Source: data.Source,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "无效的请求体",
			requestBody:    "invalid json",
			mockSetup:      func(mock *MockInterventionService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockInterventionService{}
			tt.mockSetup(mock)

			handler := NewInterventionHandler(mock)
			router := setupTestRouter()
			router.POST("/edit", handler.Edit)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/edit", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestInterventionHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockInterventionService)
		expectedStatus int
	}{
		{
			name: "成功删除",
			requestBody: domain.InterventionDeleteDTO{
				ID: 1,
			},
			mockSetup: func(mock *MockInterventionService) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "无效的请求体",
			requestBody:    "invalid json",
			mockSetup:      func(mock *MockInterventionService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockInterventionService{}
			tt.mockSetup(mock)

			handler := NewInterventionHandler(mock)
			router := setupTestRouter()
			router.POST("/delete", handler.Delete)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/delete", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
