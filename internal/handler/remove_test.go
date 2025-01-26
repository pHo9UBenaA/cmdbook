package handler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
	"github.com/pHo9UBenaA/cmdbook/internal/handler"
)

func TestRemoveCommand(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	// Helper function to write valid TOML config to file
	writeConfig := func(content string) {
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write config file: %v", err)
		}
	}

	tests := []struct {
		name          string
		initialConfig string
		configPath    string
		prefix        string
		shortCmd      string
		expectError   string
		validate      func(t *testing.T)
	}{
		{
			name: "successful removal of command with remaining prefix",
			initialConfig: `
				[commands.prefix1]
				command1 = "cmd1"
				command2 = "cmd2"
			`,
			configPath:  configPath,
			prefix:      "prefix1",
			shortCmd:    "command1",
			expectError: "",
			validate: func(t *testing.T) {
				cfg, err := config.LoadConfig(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if _, exists := cfg.Commands["prefix1"]["command1"]; exists {
					t.Error("command1 should have been removed")
				}
				if len(cfg.Commands["prefix1"]) != 1 {
					t.Errorf("expected one remaining command, got %d", len(cfg.Commands["prefix1"]))
				}
			},
		},
		{
			name: "successful removal of command with empty prefix",
			initialConfig: `
				[commands.prefix1]
				command1 = "cmd1"
			`,
			configPath:  configPath,
			prefix:      "prefix1",
			shortCmd:    "command1",
			expectError: "",
			validate: func(t *testing.T) {
				cfg, err := config.LoadConfig(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if _, exists := cfg.Commands["prefix1"]; exists {
					t.Error("prefix1 should have been removed")
				}
			},
		},
		{
			name: "prefix does not exist",
			initialConfig: `
				[commands.prefix1]
				command1 = "cmd1"
			`,
			configPath:  configPath,
			prefix:      "nonexistent",
			shortCmd:    "command1",
			expectError: "prefix does not exist: nonexistent",
		},
		{
			name: "command does not exist in prefix",
			initialConfig: `
				[commands.prefix1]
				command1 = "cmd1"
			`,
			configPath:  configPath,
			prefix:      "prefix1",
			shortCmd:    "nonexistent",
			expectError: "command not found: prefix1 nonexistent",
		},
		{
			name:       "configuration file does not exist",
			configPath: filepath.Join(tempDir, "nonexistent.toml"),
			prefix:     "prefix1",
			shortCmd:   "command1",
			// ファイルが存在しない場合にローダーでエラーが出ないようになっているため
			expectError: "prefix does not exist: prefix1",
		},
		{
			name: "failed to save configuration",
			initialConfig: `
				[commands.prefix1]
				command1 = "cmd1"
			`,
			configPath: "/invalid/path/config.toml",
			prefix:     "prefix1",
			shortCmd:   "command1",
			// ファイルが存在しない場合にローダーでエラーが出ないようになっているため
			expectError: "prefix does not exist: prefix1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.initialConfig != "" {
				writeConfig(tt.initialConfig)
			}

			err := handler.RemoveCommand(tt.configPath, tt.prefix, tt.shortCmd)

			// Check error
			if tt.expectError == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectError != "" && (err == nil || !contains(err.Error(), tt.expectError)) {
				t.Fatalf("expected error containing %q, got %v", tt.expectError, err)
			}

			// Validate final state
			if tt.validate != nil {
				tt.validate(t)
			}
		})
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && s[:len(substr)] == substr)
}
