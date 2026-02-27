package domain

import (
	"strings"
	"testing"
)

func TestGenerateSID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"生成 SID 格式正确"},
		{"生成 SID 唯一性"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sid := GenerateSID()

			// 验证格式: LB-{uuid}
			if !strings.HasPrefix(sid, "LB-") {
				t.Errorf("GenerateSID() = %v, 期望以 'LB-' 开头", sid)
			}

			// 验证长度 (LB- + 36位UUID)
			if len(sid) != 39 {
				t.Errorf("GenerateSID() 长度 = %d, 期望 39", len(sid))
			}
		})
	}

	// 测试唯一性
	t.Run("生成的 SID 应该唯一", func(t *testing.T) {
		sid1 := GenerateSID()
		sid2 := GenerateSID()

		if sid1 == sid2 {
			t.Errorf("GenerateSID() 生成了重复的 SID: %v", sid1)
		}
	})
}

func TestCalculateStep(t *testing.T) {
	tests := []struct {
		name       string
		pid        string
		parentStep int
		want       int
	}{
		{
			name:       "根标签 (pid=zhinaolabel)",
			pid:        "zhinaolabel",
			parentStep: 0,
			want:       1,
		},
		{
			name:       "二级标签",
			pid:        "LB-parent-id",
			parentStep: 1,
			want:       2,
		},
		{
			name:       "三级标签",
			pid:        "LB-parent-id",
			parentStep: 2,
			want:       3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateStep(tt.pid, tt.parentStep)
			if got != tt.want {
				t.Errorf("CalculateStep(%v, %v) = %v, 期望 %v", tt.pid, tt.parentStep, got, tt.want)
			}
		})
	}
}

func TestLabel_TableName(t *testing.T) {
	label := Label{}
	want := "relate_list"
	got := label.TableName()

	if got != want {
		t.Errorf("TableName() = %v, 期望 %v", got, want)
	}
}
