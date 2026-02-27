package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/charley/riskshield/internal/domain"
	"gorm.io/gorm"
)

func TestRiskEngineRepository_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		filters   *domain.RiskEngineQueryDTO
		mockSetup func(sqlmock.Sqlmock)
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "查询所有记录",
			filters: &domain.RiskEngineQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock COUNT query
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `riskengine`")).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				// Mock SELECT query
				rows := sqlmock.NewRows([]string{"id", "tenantID", "app", "location", "firstLabel", "secondLabel", "thirdLabel", "queryDeal", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "app1", "location1", "first1", "second1", "third1", 0, 1, "user1", "user1", time.Now(), time.Now()).
					AddRow(2, "tenant1", "app2", "location2", "first2", "second2", "third2", 1, 2, "user2", "user2", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `riskengine` LIMIT ?")).
					WithArgs(10).
					WillReturnRows(rows)
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "按 thirdLabel 过滤",
			filters: &domain.RiskEngineQueryDTO{
				Page:       1,
				Size:       10,
				ThirdLabel: stringPtr("test_label"),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `riskengine` WHERE thirdLabel = ?")).
					WithArgs("test_label").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "app", "location", "firstLabel", "secondLabel", "thirdLabel", "queryDeal", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "app1", "location1", "first1", "second1", "test_label", 0, 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `riskengine` WHERE thirdLabel = ? LIMIT ?")).
					WithArgs("test_label", 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "按 locale 过滤",
			filters: &domain.RiskEngineQueryDTO{
				Page:   1,
				Size:   10,
				Locale: stringPtr("zh-CN"),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `riskengine` WHERE location = ?")).
					WithArgs("zh-CN").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "app", "location", "firstLabel", "secondLabel", "thirdLabel", "queryDeal", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "app1", "zh-CN", "first1", "second1", "third1", 0, 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `riskengine` WHERE location = ? LIMIT ?")).
					WithArgs("zh-CN", 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.mockSetup(mock)

			repo := NewRiskEngineRepository(db)
			engines, total, err := repo.FindAll(context.Background(), tt.filters)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(engines) != tt.wantCount {
					t.Errorf("FindAll() got %d engines, want %d", len(engines), tt.wantCount)
				}
				if total != tt.wantTotal {
					t.Errorf("FindAll() got total %d, want %d", total, tt.wantTotal)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestRiskEngineRepository_FindByID(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func(sqlmock.Sqlmock)
		wantNil   bool
		wantErr   bool
	}{
		{
			name: "找到记录",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "tenantID", "app", "location", "firstLabel", "secondLabel", "thirdLabel", "queryDeal", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "app1", "location1", "first1", "second1", "third1", 0, 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `riskengine` WHERE id = ? ORDER BY `riskengine`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "记录不存在",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `riskengine` WHERE id = ? ORDER BY `riskengine`.`id` LIMIT ?")).
					WithArgs(999, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantNil: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.mockSetup(mock)

			repo := NewRiskEngineRepository(db)
			engine, err := repo.FindByID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNil && engine != nil {
				t.Errorf("FindByID() expected nil, got %+v", engine)
			}

			if !tt.wantNil && engine == nil {
				t.Error("FindByID() expected non-nil, got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestRiskEngineRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		data      *domain.RiskEngineAddDTO
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "成功创建",
			data: &domain.RiskEngineAddDTO{
				App:         "test_app",
				Location:    "test_location",
				FirstLabel:  "first",
				SecondLabel: "second",
				ThirdLabel:  "third",
				QueryDeal:   0,
				Source:      1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `riskengine`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.mockSetup(mock)

			repo := NewRiskEngineRepository(db)
			engine, err := repo.Create(context.Background(), tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && engine == nil {
				t.Error("Create() expected non-nil engine")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestRiskEngineRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		data      *domain.RiskEngineEditDTO
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "成功更新",
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
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Check if exists
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `riskengine` WHERE id = ? ORDER BY `riskengine`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(rows)

				// Update
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `riskengine`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "记录不存在",
			id:   999,
			data: &domain.RiskEngineEditDTO{
				ID:          999,
				App:         "updated_app",
				Location:    "updated_location",
				FirstLabel:  "updated_first",
				SecondLabel: "updated_second",
				ThirdLabel:  "updated_third",
				QueryDeal:   1,
				Source:      2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `riskengine` WHERE id = ? ORDER BY `riskengine`.`id` LIMIT ?")).
					WithArgs(999, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.mockSetup(mock)

			repo := NewRiskEngineRepository(db)
			engine, err := repo.Update(context.Background(), tt.id, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && engine == nil {
				t.Error("Update() expected non-nil engine")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestRiskEngineRepository_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "成功删除",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `riskengine` WHERE id = ?")).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "记录不存在",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `riskengine` WHERE id = ?")).
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.mockSetup(mock)

			repo := NewRiskEngineRepository(db)
			err := repo.Delete(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
