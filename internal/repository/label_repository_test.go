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

func TestLabelRepository_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		filters   *domain.LabelQueryDTO
		mockSetup func(sqlmock.Sqlmock)
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "查询所有记录",
			filters: &domain.LabelQueryDTO{
				Page: 1,
				Size: 10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock COUNT query
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `relate_list`")).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				// Mock SELECT query
				rows := sqlmock.NewRows([]string{"id", "tenantID", "class", "sid", "pid", "title", "value", "step", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "1", "zhinaolabel", "LB-001", "zhinaolabel", "标签1", "LB-001", 1, 2, "user1", "user1", time.Now(), time.Now()).
					AddRow(2, "1", "zhinaolabel", "LB-002", "LB-001", "标签2", "LB-002", 2, 2, "user2", "user2", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` LIMIT ?")).
					WithArgs(10).
					WillReturnRows(rows)
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "按 PID 过滤",
			filters: &domain.LabelQueryDTO{
				Page: 1,
				Size: 10,
				PID:  "zhinaolabel",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `relate_list` WHERE pid = ?")).
					WithArgs("zhinaolabel").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "class", "sid", "pid", "title", "value", "step", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "1", "zhinaolabel", "LB-001", "zhinaolabel", "标签1", "LB-001", 1, 2, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE pid = ? LIMIT ?")).
					WithArgs("zhinaolabel", 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "按 Title 模糊查询",
			filters: &domain.LabelQueryDTO{
				Page:  1,
				Size:  10,
				Title: "测试",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `relate_list` WHERE title LIKE ?")).
					WithArgs("%测试%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "class", "sid", "pid", "title", "value", "step", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "1", "zhinaolabel", "LB-001", "zhinaolabel", "测试标签", "LB-001", 1, 2, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE title LIKE ? LIMIT ?")).
					WithArgs("%测试%", 10).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "按 SID 查询",
			filters: &domain.LabelQueryDTO{
				Page: 1,
				Size: 10,
				SID:  "LB-001",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `relate_list` WHERE sid = ?")).
					WithArgs("LB-001").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "tenantID", "class", "sid", "pid", "title", "value", "step", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "1", "zhinaolabel", "LB-001", "zhinaolabel", "标签1", "LB-001", 1, 2, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE sid = ? LIMIT ?")).
					WithArgs("LB-001", 10).
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

			repo := NewLabelRepository(db)
			labels, total, err := repo.FindAll(context.Background(), tt.filters)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(labels) != tt.wantCount {
					t.Errorf("FindAll() got %d labels, want %d", len(labels), tt.wantCount)
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

func TestLabelRepository_FindByID(t *testing.T) {
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
				rows := sqlmock.NewRows([]string{"id", "tenantID", "class", "sid", "pid", "title", "value", "step", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "1", "zhinaolabel", "LB-001", "zhinaolabel", "标签1", "LB-001", 1, 2, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE id = ? ORDER BY `relate_list`.`id` LIMIT ?")).
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE id = ? ORDER BY `relate_list`.`id` LIMIT ?")).
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

			repo := NewLabelRepository(db)
			label, err := repo.FindByID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNil && label != nil {
				t.Errorf("FindByID() expected nil, got %+v", label)
			}

			if !tt.wantNil && label == nil {
				t.Error("FindByID() expected non-nil, got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestLabelRepository_FindBySID(t *testing.T) {
	tests := []struct {
		name      string
		sid       string
		mockSetup func(sqlmock.Sqlmock)
		wantNil   bool
		wantErr   bool
	}{
		{
			name: "找到记录",
			sid:  "LB-001",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "tenantID", "class", "sid", "pid", "title", "value", "step", "source", "createby", "updateby", "created_at", "updated_at"}).
					AddRow(1, "1", "zhinaolabel", "LB-001", "zhinaolabel", "标签1", "LB-001", 1, 2, "user1", "user1", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE sid = ? ORDER BY `relate_list`.`id` LIMIT ?")).
					WithArgs("LB-001", 1).
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "记录不存在",
			sid:  "LB-999",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE sid = ? ORDER BY `relate_list`.`id` LIMIT ?")).
					WithArgs("LB-999", 1).
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

			repo := NewLabelRepository(db)
			label, err := repo.FindBySID(context.Background(), tt.sid)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindBySID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNil && label != nil {
				t.Errorf("FindBySID() expected nil, got %+v", label)
			}

			if !tt.wantNil && label == nil {
				t.Error("FindBySID() expected non-nil, got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestLabelRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		data      *domain.LabelAddDTO
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "成功创建",
			data: &domain.LabelAddDTO{
				PID:   "zhinaolabel",
				Title: "测试标签",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `relate_list`")).
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

			repo := NewLabelRepository(db)
			label, err := repo.Create(context.Background(), tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && label == nil {
				t.Error("Create() expected non-nil label")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestLabelRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		data      *domain.LabelEditDTO
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "成功更新",
			id:   1,
			data: &domain.LabelEditDTO{
				ID:    1,
				PID:   "zhinaolabel",
				Title: "更新后的标签",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Check if exists
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE id = ? ORDER BY `relate_list`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(rows)

				// Update
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `relate_list`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "记录不存在",
			id:   999,
			data: &domain.LabelEditDTO{
				ID:    999,
				PID:   "zhinaolabel",
				Title: "更新后的标签",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `relate_list` WHERE id = ? ORDER BY `relate_list`.`id` LIMIT ?")).
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

			repo := NewLabelRepository(db)
			label, err := repo.Update(context.Background(), tt.id, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && label == nil {
				t.Error("Update() expected non-nil label")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestLabelRepository_Delete(t *testing.T) {
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
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `relate_list` WHERE id = ?")).
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
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `relate_list` WHERE id = ?")).
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

			repo := NewLabelRepository(db)
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
