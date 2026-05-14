package domain

import (
	"strings"
	"testing"
)

func TestCounter_Count(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		mode     SortMode
		expected ResultSet
		wantErr  bool
	}{
		{
			name:  "basic count with freq sort",
			input: "Алёна\nМиша\nАлёна\nДима\n",
			mode:  SortByFrequencyDesc,
			expected: ResultSet{
				{Name: "Алёна", Count: 2},
				{Name: "Дима", Count: 1},
				{Name: "Миша", Count: 1},
			},
			wantErr: false,
		},
		{
			name:  "sort by name asc",
			input: "Алёна\nМиша\nАлёна\nДима\n",
			mode:  SortByNameAsc,
			expected: ResultSet{
				{Name: "Алёна", Count: 2},
				{Name: "Дима", Count: 1},
				{Name: "Миша", Count: 1},
			},
			wantErr: false,
		},
		{
			name:  "empty lines should be skipped",
			input: "Алёна\n\nМиша\n\nАлёна\n",
			mode:  SortByFrequencyDesc,
			expected: ResultSet{
				{Name: "Алёна", Count: 2},
				{Name: "Миша", Count: 1},
			},
			wantErr: false,
		},
		{
			name:  "single name",
			input: "Алёна\n",
			mode:  SortByFrequencyDesc,
			expected: ResultSet{
				{Name: "Алёна", Count: 1},
			},
			wantErr: false,
		},
		{
			name:  "all names same",
			input: "Алёна\nАлёна\nАлёна\n",
			mode:  SortByFrequencyDesc,
			expected: ResultSet{
				{Name: "Алёна", Count: 3},
			},
			wantErr: false,
		},
		{
			name:  "names with spaces",
			input: "  Алёна  \n  Миша  \nАлёна\n",
			mode:  SortByFrequencyDesc,
			expected: ResultSet{
				{Name: "Алёна", Count: 2},
				{Name: "Миша", Count: 1},
			},
			wantErr: false,
		},
		{
			name:     "empty file",
			input:    "",
			mode:     SortByFrequencyDesc,
			expected: ResultSet{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewCounter()
			r := strings.NewReader(tt.input)
			results, err := counter.Count(r, tt.mode)

			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(results) != len(tt.expected) {
				t.Errorf("Count() got %d results, want %d", len(results), len(tt.expected))
				return
			}

			for i := range results {
				if results[i].Name != tt.expected[i].Name || results[i].Count != tt.expected[i].Count {
					t.Errorf("Count() result[%d] = {%s, %d}, want {%s, %d}",
						i, results[i].Name, results[i].Count,
						tt.expected[i].Name, tt.expected[i].Count)
				}
			}
		})
	}
}

func TestCounter_Count_SortOrder(t *testing.T) {
	counter := NewCounter()

	input := "B\nA\nC\nB\nA\nB\n"
	r := strings.NewReader(input)
	results, err := counter.Count(r, SortByFrequencyDesc)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	if results[0].Name != "B" || results[0].Count != 3 {
		t.Errorf("First should be B:3, got %s:%d", results[0].Name, results[0].Count)
	}
	if results[1].Name != "A" || results[1].Count != 2 {
		t.Errorf("Second should be A:2, got %s:%d", results[1].Name, results[1].Count)
	}
	if results[2].Name != "C" || results[2].Count != 1 {
		t.Errorf("Third should be C:1, got %s:%d", results[2].Name, results[2].Count)
	}
}

func TestCounter_Count_DefaultSort(t *testing.T) {
	counter := NewCounter()
	input := "B\nA\nC\nB\nA\nB\n"
	r := strings.NewReader(input)

	results, err := counter.Count(r, SortMode(999))
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	if results[0].Name != "B" || results[0].Count != 3 {
		t.Errorf("First should be B:3, got %s:%d", results[0].Name, results[0].Count)
	}
}

func TestCounter_Count_LargeLine(t *testing.T) {
	counter := NewCounter()
	longName := strings.Repeat("a", 70*1024)
	input := longName + "\n" + longName + "\n"

	r := strings.NewReader(input)
	results, err := counter.Count(r, SortByFrequencyDesc)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0].Count != 2 {
		t.Errorf("Expected count 2, got %d", results[0].Count)
	}
}
