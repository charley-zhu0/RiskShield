package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/charley/riskshield/internal/domain"
	"github.com/gin-gonic/gin"
)

// MockRiskEngineService 模拟 Service
type MockRiskEngineService struct {
	QueryFunc  func(ctx context.Context, filters *domain.RiskEngineQueryDTO) (*domain.RiskEngineQueryResponse, error)
	AddFunc    func(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error)
	EditFunc   func(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error)
	DeleteFunc func(ctx context.Context, id uint) error
}

func (m *MockRiskEngineService) Query(ctx context.Context, filters *domain.RiskEngineQueryDTO) (*domain.RiskEngineQueryResponse, error) {
	return m.QueryFunc(ctx, filters)
}

func (m *MockRiskEngineService) Add(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error) {
	return m.AddFunc(ctx, data)
}

func (m *MockRiskEngineService) Edit(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error) {
	return m.EditFunc(ctx, id, data)
}

func (m *MockRiskEngineService) Delete(ctx context.Context, id uint) error {
	return m.DeleteFunc(ctx, id)
}

func TestRiskEngineHandler_Query(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockRiskEngineService)
		expectedStatus int
		expectedErrNo  int
	}{
		{
			name:        "成功查询",
			queryParams: "?page=1&size=10",
			mockSetup: func(mock *MockRiskEngineService) {
				mock.QueryFunc = func(ctx context.Context, filters *domain.RiskEngineQueryDTO) (*domain.RiskEngineQueryResponse, error) {
					return &domain.RiskEngineQueryResponse{
						Data: []domain.RiskEngine{
							{ID: 1, App: "app1", Location: "loc1", FirstLabel: "first1", SecondLabel: "second1", ThirdLabel: "third1", QueryDeal: 0, Source: 1},
						},
						Total: 1,
						Page:  1,
						Size:  10,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedErrNo:  0,
		},
		{
			name:           "参数错误",
			queryParams:    "?page=0&size=10",
			mockSetup:      func(mock *MockRiskEngineService) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrNo:  400,
		},
		{
			name:        "服务错误",
			queryParams: "?page=1&size=10",
			mockSetup: func(mock *MockRiskEngineService) {
				mock.QueryFunc = func(ctx context.Context, filters *domain.RiskEngineQueryDTO) (*domain.RiskEngineQueryResponse, error) {
					return nil, errors.New("service error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrNo:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockRiskEngineService{}
			tt.mockSetup(mockService)

			handler := NewRiskEngineHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/risk_shield/admin/riskengine/query"+tt.queryParams, nil)

			handler.Query(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Query() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var response domain.StandardResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if response.ErrNo != tt.expectedErrNo {
				t.Errorf("Query() errno = %d, want %d", response.ErrNo, tt.expectedErrNo)
			}
		})
	}
}

func TestRiskEngineHandler_Add(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockRiskEngineService)
		expectedStatus int
		expectedErrNo  int
	}{
		{
			name: "成功添加",
			requestBody: domain.RiskEngineAddDTO{
				App:         "test_app",
				Location:    "test_location",
				FirstLabel:  "first",
				SecondLabel: "second",
				ThirdLabel:  "third",
				QueryDeal:   0,
				Source:      1,
			},
			mockSetup: func(mock *MockRiskEngineService) {
				mock.AddFunc = func(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error) {
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
			expectedStatus: http.StatusOK,
			expectedErrNo:  0,
		},
		{
			name:           "参数错误",
			requestBody:    map[string]interface{}{"app": ""},
			mockSetup:      func(mock *MockRiskEngineService) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrNo:  400,
		},
		{
			name: "服务错误",
			requestBody: domain.RiskEngineAddDTO{
				App:         "test_app",
				Location:    "test_location",
				FirstLabel:  "first",
				SecondLabel: "second",
				ThirdLabel:  "third",
				QueryDeal:   0,
				Source:      1,
			},
			mockSetup: func(mock *MockRiskEngineService) {
				mock.AddFunc = func(ctx context.Context, data *domain.RiskEngineAddDTO) (*domain.RiskEngine, error) {
					return nil, errors.New("service error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrNo:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockRiskEngineService{}
			tt.mockSetup(mockService)

			handler := NewRiskEngineHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/risk_shield/admin/admin/riskengine/add", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Add(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Add() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var response domain.StandardResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if response.ErrNo != tt.expectedErrNo {
				t.Errorf("Add() errno = %d, want %d", response.ErrNo, tt.expectedErrNo)
			}
		})
	}
}

func TestRiskEngineHandler_Edit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockRiskEngineService)
		expectedStatus int
		expectedErrNo  int
	}{
		{
			name: "成功编辑",
			requestBody: domain.RiskEngineEditDTO{
				ID:          1,
				App:         "updated_app",
				Location:    "updated_location",
				FirstLabel:  "updated_first",
				SecondLabel: "updated_second",
				ThirdLabel:  "updated_third",
				QueryDeal:   1,
				Source:      2,
			},
			mockSetup: func(mock *MockRiskEngineService) {
				mock.EditFunc = func(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error) {
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
			expectedStatus: http.StatusOK,
			expectedErrNo:  0,
		},
		{
			name:           "参数错误",
			requestBody:    map[string]interface{}{"id": 0},
			mockSetup:      func(mock *MockRiskEngineService) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrNo:  400,
		},
		{
			name: "服务错误",
			requestBody: domain.RiskEngineEditDTO{
				ID:          1,
				App:         "updated_app",
				Location:    "updated_location",
				FirstLabel:  "updated_first",
				SecondLabel: "updated_second",
				ThirdLabel:  "updated_third",
				QueryDeal:   1,
				Source:      2,
			},
			mockSetup: func(mock *MockRiskEngineService) {
				mock.EditFunc = func(ctx context.Context, id uint, data *domain.RiskEngineEditDTO) (*domain.RiskEngine, error) {
					return nil, errors.New("service error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrNo:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockRiskEngineService{}
			tt.mockSetup(mockService)

			handler := NewRiskEngineHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/risk_shield/admin/admin/riskengine/edit", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Edit(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Edit() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var response domain.StandardResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if response.ErrNo != tt.expectedErrNo {
				t.Errorf("Edit() errno = %d, want %d", response.ErrNo, tt.expectedErrNo)
			}
		})
	}
}

func TestRiskEngineHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockRiskEngineService)
		expectedStatus int
		expectedErrNo  int
	}{
		{
			name:        "成功删除",
			requestBody: domain.RiskEngineDeleteDTO{ID: 1},
			mockSetup: func(mock *MockRiskEngineService) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedErrNo:  0,
		},
		{
			name:           "参数错误",
			requestBody:    map[string]interface{}{"id": 0},
			mockSetup:      func(mock *MockRiskEngineService) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrNo:  400,
		},
		{
			name:        "服务错误",
			requestBody: domain.RiskEngineDeleteDTO{ID: 1},
			mockSetup: func(mock *MockRiskEngineService) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return errors.New("service error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrNo:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockRiskEngineService{}
			tt.mockSetup(mockService)

			handler := NewRiskEngineHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/risk_shield/admin/admin/riskengine/delete", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Delete(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Delete() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var response domain.StandardResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if response.ErrNo != tt.expectedErrNo {
				t.Errorf("Delete() errno = %d, want %d", response.ErrNo, tt.expectedErrNo)
			}
		})
	}
}
