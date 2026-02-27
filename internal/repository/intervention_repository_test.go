package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/charley/riskshield/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() {
		sqlDB.Close()
	}

	return gormDB, mock, cleanup
}

func TestInterventionRepository_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		filters   *domain.InterventionQueryDTO
		mockSetup func(sqlmock.Sqlmock)
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "查询所有记录",
			filters: &domain.InterventionQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock COUNT query
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `intervention`")).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				// Mock SELECT query
				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "test query 1", "test answer 1", "hash1", 1, "user1", "user1", time.Now(), time.Now()).
					AddRow(2, "tenant1", "test query 2", "test answer 2", "hash2", 2, "user2", "user2", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` LIMIT ?")).
					WithArgs(10).
					WillReturnRows(rows)
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "精确查询",
			filters: &domain.InterventionQueryDTO{
				Page:  1,
				Size:  10,
				Query: "test query",
				Fuzzy: false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				queryHash := domain.CalculateQueryHash("test query")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `intervention` WHERE queryHash = ?")).
					WithArgs(queryHash).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "test query", "test answer", queryHash, 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE queryHash = ? LIMIT ?")).
					WithArgs(queryHash, 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "模糊查询",
			filters: &domain.InterventionQueryDTO{
				Page:  1,
				Size:  10,
				Query: "test",
				Fuzzy: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `intervention` WHERE query LIKE ?")).
					WithArgs("%test%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "test query 1", "test answer 1", "hash1", 1, "user1", "user1", time.Now(), time.Now()).
					AddRow(2, "tenant1", "test query 2", "test answer 2", "hash2", 2, "user2", "user2", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE query LIKE ? LIMIT ?")).
					WithArgs("%test%", 10).
					WillReturnRows(rows)
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "模糊查询 - 转义反斜杠",
			filters: &domain.InterventionQueryDTO{
				Page:  1,
				Size:  10,
				Query: `test\query`,
				Fuzzy: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `intervention` WHERE query LIKE ?")).
					WithArgs(`%test\\query%`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", `test\query`, "test answer", "hash1", 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE query LIKE ? LIMIT ?")).
					WithArgs(`%test\\query%`, 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "模糊查询 - 转义百分号",
			filters: &domain.InterventionQueryDTO{
				Page:  1,
				Size:  10,
				Query: "test%query",
				Fuzzy: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `intervention` WHERE query LIKE ?")).
					WithArgs(`%test\%query%`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "test%query", "test answer", "hash1", 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE query LIKE ? LIMIT ?")).
					WithArgs(`%test\%query%`, 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "模糊查询 - 转义下划线",
			filters: &domain.InterventionQueryDTO{
				Page:  1,
				Size:  10,
				Query: "test_query",
				Fuzzy: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `intervention` WHERE query LIKE ?")).
					WithArgs(`%test\_query%`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "test_query", "test answer", "hash1", 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE query LIKE ? LIMIT ?")).
					WithArgs(`%test\_query%`, 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "模糊查询 - 转义所有特殊字符",
			filters: &domain.InterventionQueryDTO{
				Page:  1,
				Size:  10,
				Query: `test\%_query`,
				Fuzzy: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `intervention` WHERE query LIKE ?")).
					WithArgs(`%test\\\%\_query%`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", `test\%_query`, "test answer", "hash1", 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE query LIKE ? LIMIT ?")).
					WithArgs(`%test\\\%\_query%`, 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "按 source 过滤",
			filters: &domain.InterventionQueryDTO{
				Page:   1,
				Size:   10,
				Source: intPtr(2),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `intervention` WHERE source = ?")).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(2, "tenant1", "test query 2", "test answer 2", "hash2", 2, "user2", "user2", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE source = ? LIMIT ?")).
					WithArgs(2, 10).
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

			repo := NewInterventionRepository(db)
			interventions, total, err := repo.FindAll(context.Background(), tt.filters)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(interventions) != tt.wantCount {
					t.Errorf("FindAll() got %d interventions, want %d", len(interventions), tt.wantCount)
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

func TestInterventionRepository_FindByID(t *testing.T) {
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
				rows := sqlmock.NewRows([]string{"id", "tenantID", "query", "answer", "queryHash", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "tenant1", "test query", "test answer", "hash1", 1, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE id = ? ORDER BY `intervention`.`id` LIMIT ?")).
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE id = ? ORDER BY `intervention`.`id` LIMIT ?")).
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

			repo := NewInterventionRepository(db)
			intervention, err := repo.FindByID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNil && intervention != nil {
				t.Errorf("FindByID() expected nil, got %+v", intervention)
			}

			if !tt.wantNil && intervention == nil {
				t.Error("FindByID() expected non-nil, got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestInterventionRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		data      *domain.InterventionAddDTO
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "成功创建",
			data: &domain.InterventionAddDTO{
				Query:  "test query",
				Answer: "test answer",
				Source: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `intervention`")).
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

			repo := NewInterventionRepository(db)
			intervention, err := repo.Create(context.Background(), tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && intervention == nil {
				t.Error("Create() expected non-nil intervention")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestInterventionRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		data      *domain.InterventionEditDTO
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "成功更新",
			id:   1,
			data: &domain.InterventionEditDTO{
				ID:     1,
				Query:  "updated query",
				Answer: "updated answer",
				Source: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Check if exists
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE id = ? ORDER BY `intervention`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(rows)

				// Update
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `intervention`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
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
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `intervention` WHERE id = ? ORDER BY `intervention`.`id` LIMIT ?")).
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

			repo := NewInterventionRepository(db)
			intervention, err := repo.Update(context.Background(), tt.id, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && intervention == nil {
				t.Error("Update() expected non-nil intervention")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestInterventionRepository_Delete(t *testing.T) {
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
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `intervention` WHERE id = ?")).
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
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `intervention` WHERE id = ?")).
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

			repo := NewInterventionRepository(db)
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

func intPtr(i int) *int {
	return &i
}
