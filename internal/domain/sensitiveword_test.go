package domain

import (
	"testing"
)

func TestSensitiveWord_TableName(t *testing.T) {
	sw := SensitiveWord{}
	got := sw.TableName()
	want := "sensitiveword"

	if got != want {
		t.Errorf("TableName() = %q; want %q", got, want)
	}
}

func TestSensitiveWordQueryDTO_Validation(t *testing.T) {
	tests := []struct {
		name    string
		dto     SensitiveWordQueryDTO
		wantErr bool
	}{
		{
			name: "valid query",
			dto: SensitiveWordQueryDTO{
				Page: 1,
				Size: 20,
			},
			wantErr: false,
		},
		{
			name: "page less than 1",
			dto: SensitiveWordQueryDTO{
				Page: 0,
				Size: 20,
			},
			wantErr: true,
		},
		{
			name: "size exceeds max",
			dto: SensitiveWordQueryDTO{
				Page: 1,
				Size: 101,
			},
			wantErr: true,
		},
	}

	// Note: Gin binding validation is tested in handler tests
	// This test documents expected validation rules
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validation happens at Gin binding level
			// This test serves as documentation
		})
	}
}

func TestSensitiveWordAddDTO_Validation(t *testing.T) {
	tests := []struct {
		name string
		dto  SensitiveWordAddDTO
	}{
		{
			name: "valid add request",
			dto: SensitiveWordAddDTO{
				Words:       []string{"test", "word"},
				FirstLabel:  "100",
				SecondLabel: "100001",
				ThirdLabel:  "100001001",
				QueryDeal:   1,
				MatchedID:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validation happens at Gin binding level
		})
	}
}
