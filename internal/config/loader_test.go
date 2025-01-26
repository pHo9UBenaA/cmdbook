package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		wantCmds map[string]map[string]string
		wantErr  bool
	}{
		{
			name: "file does not exist",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.toml")
			},
			wantCmds: map[string]map[string]string{},
			wantErr:  false,
		},
		{
			name: "invalid toml content",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				path := filepath.Join(dir, "invalid.toml")
				if err := os.WriteFile(path, []byte("invalid toml content"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: true,
		},
		{
			name: "empty file",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				path := filepath.Join(dir, "empty.toml")
				if err := os.WriteFile(path, []byte{}, 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantCmds: map[string]map[string]string{},
			wantErr:  false,
		},
		{
			name: "valid toml content",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				path := filepath.Join(dir, "valid.toml")
				content := `
				[commands]
				[commands.test]
				action = "echo hello"
				`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantCmds: map[string]map[string]string{
				"test": {"action": "echo hello"},
			},
			wantErr: false,
		},
		{
			name: "read error (path is directory)",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				path := filepath.Join(dir, "directory")
				if err := os.Mkdir(path, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)
			cfg, err := config.LoadConfig(path)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(cfg.Commands) != len(tt.wantCmds) {
					t.Errorf("Commands length mismatch: got %d, want %d", len(cfg.Commands), len(tt.wantCmds))
				}
				for key, expected := range tt.wantCmds {
					actual, ok := cfg.Commands[key]
					if !ok {
						t.Errorf("missing command key: %s", key)
						continue
					}
					for k, v := range expected {
						if actual[k] != v {
							t.Errorf("command %s: got %s = %s, want %s", key, k, actual[k], v)
						}
					}
				}
			}
		})
	}
}
