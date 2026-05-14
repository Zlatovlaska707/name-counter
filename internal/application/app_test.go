package application

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"name-counter/internal/domain"
)

func TestApp_Run(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		sortMode   domain.SortMode
		expected   string
		wantErr    bool
		errMessage string
	}{
		{
			name:     "basic counting with freq sort",
			content:  "Алёна\nМиша\nАлёна\nДима\n",
			sortMode: domain.SortByFrequencyDesc,
			expected: "Алёна: 2\nДима: 1\nМиша: 1\n",
			wantErr:  false,
		},
		{
			name:     "sort by name",
			content:  "Миша\nАлёна\nДима\nАлёна\n",
			sortMode: domain.SortByNameAsc,
			expected: "Алёна: 2\nДима: 1\nМиша: 1\n",
			wantErr:  false,
		},
		{
			name:     "empty lines",
			content:  "Алёна\n\nМиша\n\nАлёна\n",
			sortMode: domain.SortByFrequencyDesc,
			expected: "Алёна: 2\nМиша: 1\n",
			wantErr:  false,
		},
		{
			name:     "single name",
			content:  "Алёна\n",
			sortMode: domain.SortByFrequencyDesc,
			expected: "Алёна: 1\n",
			wantErr:  false,
		},
		{
			name:     "empty file",
			content:  "",
			sortMode: domain.SortByFrequencyDesc,
			expected: "",
			wantErr:  false,
		},
		{
			name:       "file not found",
			content:    "",
			sortMode:   domain.SortByFrequencyDesc,
			expected:   "",
			wantErr:    true,
			errMessage: "не удалось открыть файл",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string

			if tt.name == "file not found" {
				filePath = "/nonexistent/file.txt"
			} else {
				tmpFile, err := os.CreateTemp("", "test-*.txt")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(tmpFile.Name())
				filePath = tmpFile.Name()

				if _, err := tmpFile.WriteString(tt.content); err != nil {
					t.Fatalf("Failed to write to temp file: %v", err)
				}
				if err := tmpFile.Close(); err != nil {
					t.Fatalf("Failed to close temp file: %v", err)
				}
			}

			app := New(&domain.Config{
				FilePath: filePath,
				SortMode: tt.sortMode,
			})

			var out bytes.Buffer
			var errOut bytes.Buffer

			err := app.Run(&out, &errOut)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Run() expected error, got nil")
				} else if tt.errMessage != "" && !strings.Contains(err.Error(), tt.errMessage) {
					t.Errorf("Run() error = %v, should contain %v", err, tt.errMessage)
				}
				return
			}

			if err != nil {
				t.Errorf("Run() unexpected error: %v", err)
				return
			}

			if out.String() != tt.expected {
				t.Errorf("Run() output = %q, want %q", out.String(), tt.expected)
			}
		})
	}
}

func TestApp_RunWithDefaultOutput(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "Алёна\nМиша\nАлёна\n"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	app := New(&domain.Config{
		FilePath: tmpFile.Name(),
		SortMode: domain.SortByFrequencyDesc,
	})

	err = app.RunWithDefaultOutput()
	if err != nil {
		t.Errorf("RunWithDefaultOutput() unexpected error: %v", err)
	}
}

func TestApp_Run_FileError(t *testing.T) {
	app := New(&domain.Config{
		FilePath: "/tmp/nonexistent_file_12345.txt",
		SortMode: domain.SortByFrequencyDesc,
	})

	var out bytes.Buffer
	var errOut bytes.Buffer

	err := app.Run(&out, &errOut)
	if err == nil {
		t.Error("Run() expected error for non-existent file, got nil")
	}
}

func TestApp_Run_LargeFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "large-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := strings.Repeat("Александр\n", 5000) + strings.Repeat("Мария\n", 3000) + strings.Repeat("Иван\n", 2000)
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	app := New(&domain.Config{
		FilePath: tmpFile.Name(),
		SortMode: domain.SortByFrequencyDesc,
	})

	var out bytes.Buffer
	var errOut bytes.Buffer

	err = app.Run(&out, &errOut)
	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Александр: 5000") {
		t.Error("Expected Александр count 5000 in output")
	}
	if !strings.Contains(output, "Мария: 3000") {
		t.Error("Expected Мария count 3000 in output")
	}
	if !strings.Contains(output, "Иван: 2000") {
		t.Error("Expected Иван count 2000 in output")
	}
}

func TestApp_Run_UnicodeNames(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "unicode-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "José\nMüller\nFrançois\nJosé\nŁukasz\nMüller\n"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	app := New(&domain.Config{
		FilePath: tmpFile.Name(),
		SortMode: domain.SortByNameAsc,
	})

	var out bytes.Buffer
	var errOut bytes.Buffer

	err = app.Run(&out, &errOut)
	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}

	expected := "François: 1\nJosé: 2\nMüller: 2\nŁukasz: 1\n"
	if out.String() != expected {
		t.Errorf("Run() output = %q, want %q", out.String(), expected)
	}
}

func TestApp_Run_LongNames(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "long-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	longName := strings.Repeat("Александр", 100) // Very long name
	content := longName + "\n" + longName + "\n" + "Short\n"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	app := New(&domain.Config{
		FilePath: tmpFile.Name(),
		SortMode: domain.SortByFrequencyDesc,
	})

	var out bytes.Buffer
	var errOut bytes.Buffer

	err = app.Run(&out, &errOut)
	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), longName) {
		t.Error("Output missing long name")
	}
	if !strings.Contains(out.String(), "Short: 1") {
		t.Error("Output missing short name")
	}
}

func TestApp_Run_WriteError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "Алёна\nМиша\n"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	app := New(&domain.Config{
		FilePath: tmpFile.Name(),
		SortMode: domain.SortByFrequencyDesc,
	})

	closedWriter := &closedWriter{}
	errOut := &bytes.Buffer{}

	err = app.Run(closedWriter, errOut)
	if err == nil {
		t.Error("Expected error when writing to closed writer, got nil")
	}
}

type closedWriter struct{}

func (cw *closedWriter) Write(p []byte) (n int, err error) {
	return 0, os.ErrClosed
}
