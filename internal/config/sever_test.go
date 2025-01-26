package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		setup   func(t *testing.T) string
		wantErr bool
	}{
		{
			name: "successful save",
			cfg: &Config{
				Commands: map[string]map[string]string{
					"test": {"action": "echo hello"},
				},
			},
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "toml")
			},
			wantErr: false,
		},
		{
			name: "marshal error (nil config)",
			cfg:  nil,
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nil.toml")
			},
			wantErr: false,
		},
		{
			name: "write error (path is directory)",
			cfg: &Config{
				Commands: map[string]map[string]string{},
			},
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
			err := SaveConfig(tt.cfg, path)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.cfg != nil {
				loaded, err := LoadConfig(path)
				if err != nil {
					t.Fatalf("failed to load saved config: %v", err)
				}

				if len(loaded.Commands) != len(tt.cfg.Commands) {
					t.Errorf("saved commands count mismatch: got %d, want %d", len(loaded.Commands), len(tt.cfg.Commands))
				}
				for key, expected := range tt.cfg.Commands {
					actual, ok := loaded.Commands[key]
					if !ok {
						t.Errorf("missing command key in saved config: %s", key)
						continue
					}
					for k, v := range expected {
						if actual[k] != v {
							t.Errorf("saved command %s.%s: got %s, want %s", key, k, actual[k], v)
						}
					}
				}
			}
		})
	}
}
