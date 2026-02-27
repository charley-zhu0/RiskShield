package repository

import "testing"

// TestEscapeLikePattern 测试 LIKE 查询转义函数
func TestEscapeLikePattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "无特殊字符",
			input:    "test query",
			expected: "test query",
		},
		{
			name:     "转义反斜杠",
			input:    `test\query`,
			expected: `test\\query`,
		},
		{
			name:     "转义百分号",
			input:    "test%query",
			expected: `test\%query`,
		},
		{
			name:     "转义下划线",
			input:    "test_query",
			expected: `test\_query`,
		},
		{
			name:     "转义所有特殊字符",
			input:    `test\%_query`,
			expected: `test\\\%\_query`,
		},
		{
			name:     "多个反斜杠",
			input:    `test\\query`,
			expected: `test\\\\query`,
		},
		{
			name:     "多个百分号",
			input:    "test%%query",
			expected: `test\%\%query`,
		},
		{
			name:     "多个下划线",
			input:    "test__query",
			expected: `test\_\_query`,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有特殊字符",
			input:    `\%_`,
			expected: `\\\%\_`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapeLikeString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapeLikeString(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
