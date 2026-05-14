package domain

import (
	"strings"
	"testing"
)

func TestNameCount_String(t *testing.T) {
	tests := []struct {
		name     string
		nc       NameCount
		expected string
	}{
		{
			name:     "basic name count",
			nc:       NameCount{Name: "Алёна", Count: 5},
			expected: "Алёна: 5",
		},
		{
			name:     "zero count",
			nc:       NameCount{Name: "Миша", Count: 0},
			expected: "Миша: 0",
		},
		{
			name:     "large count",
			nc:       NameCount{Name: "Дима", Count: 1000000},
			expected: "Дима: 1000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nc.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResultSet_String(t *testing.T) {
	tests := []struct {
		name     string
		rs       ResultSet
		expected string
	}{
		{
			name:     "empty result set",
			rs:       ResultSet{},
			expected: "",
		},
		{
			name: "single result",
			rs: ResultSet{
				{Name: "Алёна", Count: 2},
			},
			expected: "Алёна: 2\n",
		},
		{
			name: "multiple results",
			rs: ResultSet{
				{Name: "Алёна", Count: 2},
				{Name: "Дима", Count: 1},
			},
			expected: "Алёна: 2\nДима: 1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rs.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResultSet_String_NewlineHandling(t *testing.T) {
	rs := ResultSet{
		{Name: "Name1", Count: 1},
		{Name: "Name2", Count: 2},
	}
	result := rs.String()

	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines from split, got %d", len(lines))
	}
}
