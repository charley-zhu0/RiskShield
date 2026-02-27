package repository

import (
	"context"
	"testing"

	"github.com/charley/riskshield/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&domain.SensitiveWord{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestSensitiveWordRepository_CreateBatch(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		data    *domain.SensitiveWordAddDTO
		wantLen int
		wantErr bool
	}{
		{
			name: "create single word",
			data: &domain.SensitiveWordAddDTO{
				Words:       []string{"test"},
				FirstLabel:  "100",
				SecondLabel: "100001",
				ThirdLabel:  "100001001",
				QueryDeal:   1,
				MatchedID:   0,
			},
			wantLen: 1,
		},
		{
			name: "create multiple words",
			data: &domain.SensitiveWordAddDTO{
				Words:       []string{"word1", "word2", "word3"},
				FirstLabel:  "100",
				SecondLabel: "100001",
				ThirdLabel:  "100001001",
				QueryDeal:   1,
				MatchedID:   0,
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.CreateBatch(ctx, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("CreateBatch() got %d records, want %d", len(got), tt.wantLen)
			}

			// 验证字段
			if !tt.wantErr && len(got) > 0 {
				if got[0].FirstLabel != tt.data.FirstLabel {
					t.Errorf("FirstLabel = %v, want %v", got[0].FirstLabel, tt.data.FirstLabel)
				}
			}
		})
	}
}

func TestSensitiveWordRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	// 准备测试数据
	testData := &domain.SensitiveWordAddDTO{
		Words:       []string{"test1", "test2", "cs"},
		FirstLabel:  "100",
		SecondLabel: "100001",
		ThirdLabel:  "100001001",
		QueryDeal:   1,
		MatchedID:   0,
	}
	_, err := repo.CreateBatch(ctx, testData)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	tests := []struct {
		name      string
		filters   *domain.SensitiveWordQueryDTO
		wantCount int64
		wantErr   bool
	}{
		{
			name: "query all",
			filters: &domain.SensitiveWordQueryDTO{
				Page: 1,
				Size: 10,
			},
			wantCount: 3,
		},
		{
			name: "query with word filter",
			filters: &domain.SensitiveWordQueryDTO{
				Page: 1,
				Size: 10,
				Word: "cs",
			},
			wantCount: 1,
		},
		{
			name: "query with pagination",
			filters: &domain.SensitiveWordQueryDTO{
				Page: 1,
				Size: 2,
			},
			wantCount: 3, // total count
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, total, err := repo.FindAll(ctx, tt.filters)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if total != tt.wantCount {
				t.Errorf("FindAll() total = %d, want %d", total, tt.wantCount)
			}

			if !tt.wantErr && got == nil {
				t.Error("FindAll() returned nil slice")
			}
		})
	}
}

func TestSensitiveWordRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	// 创建测试数据
	testData := &domain.SensitiveWordAddDTO{
		Words:       []string{"test"},
		FirstLabel:  "100",
		SecondLabel: "100001",
		QueryDeal:   1,
	}
	created, err := repo.CreateBatch(ctx, testData)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	tests := []struct {
		name    string
		id      uint
		wantNil bool
		wantErr bool
	}{
		{
			name:    "find existing",
			id:      created[0].ID,
			wantNil: false,
		},
		{
			name:    "find non-existing",
			id:      9999,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByID(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("FindByID() got nil = %v, want nil = %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestSensitiveWordRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	// 创建测试数据
	testData := &domain.SensitiveWordAddDTO{
		Words:       []string{"original"},
		FirstLabel:  "100",
		SecondLabel: "100001",
		QueryDeal:   1,
	}
	created, err := repo.CreateBatch(ctx, testData)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	tests := []struct {
		name    string
		id      uint
		data    *domain.SensitiveWordEditDTO
		wantErr bool
	}{
		{
			name: "update existing",
			id:   created[0].ID,
			data: &domain.SensitiveWordEditDTO{
				ID:          created[0].ID,
				Word:        "updated",
				FirstLabel:  "200",
				SecondLabel: "200001",
				QueryDeal:   2,
				MatchedID:   1,
			},
		},
		{
			name: "update non-existing",
			id:   9999,
			data: &domain.SensitiveWordEditDTO{
				ID:          9999,
				Word:        "test",
				FirstLabel:  "100",
				SecondLabel: "100001",
				QueryDeal:   1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Update(ctx, tt.id, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Word != tt.data.Word {
					t.Errorf("Update() Word = %v, want %v", got.Word, tt.data.Word)
				}
				if got.FirstLabel != tt.data.FirstLabel {
					t.Errorf("Update() FirstLabel = %v, want %v", got.FirstLabel, tt.data.FirstLabel)
				}
			}
		})
	}
}

func TestSensitiveWordRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSensitiveWordRepository(db)
	ctx := context.Background()

	// 创建测试数据
	testData := &domain.SensitiveWordAddDTO{
		Words:       []string{"to_delete"},
		FirstLabel:  "100",
		SecondLabel: "100001",
		QueryDeal:   1,
	}
	created, err := repo.CreateBatch(ctx, testData)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name: "delete existing",
			id:   created[0].ID,
		},
		{
			name:    "delete non-existing",
			id:      9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 验证删除成功
			if !tt.wantErr {
				found, _ := repo.FindByID(ctx, tt.id)
				if found != nil {
					t.Error("Delete() record still exists")
				}
			}
		})
	}
}
