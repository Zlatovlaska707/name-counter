package cli

import (
	"flag"
	"os"
	"testing"

	"name-counter/internal/domain"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantConfig  *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "file as positional argument",
			args: []string{"cmd", "test.txt"},
			wantConfig: &Config{
				FilePath:  "test.txt",
				SortMode:  domain.SortByFrequencyDesc,
				PprofAddr: "",
				ShowHelp:  false,
			},
			wantErr: false,
		},
		{
			name: "file with -file flag",
			args: []string{"cmd", "-file", "test.txt"},
			wantConfig: &Config{
				FilePath:  "test.txt",
				SortMode:  domain.SortByFrequencyDesc,
				PprofAddr: "",
				ShowHelp:  false,
			},
			wantErr: false,
		},
		{
			name: "sort by name",
			args: []string{"cmd", "-file", "test.txt", "-sort", "name"},
			wantConfig: &Config{
				FilePath:  "test.txt",
				SortMode:  domain.SortByNameAsc,
				PprofAddr: "",
				ShowHelp:  false,
			},
			wantErr: false,
		},
		{
			name: "sort by frequency",
			args: []string{"cmd", "-file", "test.txt", "-sort", "freq"},
			wantConfig: &Config{
				FilePath:  "test.txt",
				SortMode:  domain.SortByFrequencyDesc,
				PprofAddr: "",
				ShowHelp:  false,
			},
			wantErr: false,
		},
		{
			name: "sort by frequency alias",
			args: []string{"cmd", "-file", "test.txt", "-sort", "frequency"},
			wantConfig: &Config{
				FilePath:  "test.txt",
				SortMode:  domain.SortByFrequencyDesc,
				PprofAddr: "",
				ShowHelp:  false,
			},
			wantErr: false,
		},
		{
			name: "with pprof address",
			args: []string{"cmd", "-file", "test.txt", "-pprof", ":6060"},
			wantConfig: &Config{
				FilePath:  "test.txt",
				SortMode:  domain.SortByFrequencyDesc,
				PprofAddr: ":6060",
				ShowHelp:  false,
			},
			wantErr: false,
		},
		{
			name: "show help",
			args: []string{"cmd", "-help"},
			wantConfig: &Config{
				FilePath:  "",
				SortMode:  0,
				PprofAddr: "",
				ShowHelp:  true,
			},
			wantErr: false,
		},
		{
			name:        "no file specified",
			args:        []string{"cmd"},
			wantErr:     true,
			errContains: "не указан путь к файлу",
		},
		{
			name:        "invalid sort mode",
			args:        []string{"cmd", "-file", "test.txt", "-sort", "invalid"},
			wantErr:     true,
			errContains: "неверный -sort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			parser := NewParser()
			cfg, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error, got nil")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Parse() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if cfg.ShowHelp != tt.wantConfig.ShowHelp {
				t.Errorf("ShowHelp = %v, want %v", cfg.ShowHelp, tt.wantConfig.ShowHelp)
			}
			if cfg.FilePath != tt.wantConfig.FilePath {
				t.Errorf("FilePath = %v, want %v", cfg.FilePath, tt.wantConfig.FilePath)
			}
			if cfg.SortMode != tt.wantConfig.SortMode {
				t.Errorf("SortMode = %v, want %v", cfg.SortMode, tt.wantConfig.SortMode)
			}
			if cfg.PprofAddr != tt.wantConfig.PprofAddr {
				t.Errorf("PprofAddr = %v, want %v", cfg.PprofAddr, tt.wantConfig.PprofAddr)
			}
		})
	}
}

func TestParser_determineFilePath(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "from -file flag",
			args:     []string{"cmd", "-file", "flagfile.txt"},
			expected: "flagfile.txt",
		},
		{
			name:     "from positional arg",
			args:     []string{"cmd", "positional.txt"},
			expected: "positional.txt",
		},
		{
			name:     "flag overrides positional",
			args:     []string{"cmd", "-file", "flagfile.txt", "positional.txt"},
			expected: "flagfile.txt",
		},
		{
			name:     "no file",
			args:     []string{"cmd"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			parser := NewParser()
			parser.filePath = flag.String("file", "", "test")
			flag.Parse()

			if len(tt.args) > 2 && tt.args[1] == "-file" {
				*parser.filePath = tt.args[2]
			} else {
				*parser.filePath = ""
			}

			result := parser.determineFilePath()
			if result != tt.expected {
				t.Errorf("determineFilePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParser_parseSortMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected domain.SortMode
		wantErr  bool
	}{
		{"freq mode", "freq", domain.SortByFrequencyDesc, false},
		{"frequency mode", "frequency", domain.SortByFrequencyDesc, false},
		{"name mode", "name", domain.SortByNameAsc, false},
		{"invalid mode", "invalid", 0, true},
		{"empty mode", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			result, err := parser.parseSortMode(tt.input)

			if tt.wantErr && err == nil {
				t.Errorf("parseSortMode() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("parseSortMode() unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("parseSortMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
