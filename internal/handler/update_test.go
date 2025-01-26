package handler_test

import (
	"os"
	"testing"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
	"github.com/pHo9UBenaA/cmdbook/internal/handler"
)

func TestUpdateCommand(t *testing.T) {
	tests := []struct {
		name          string
		initialConfig *config.Config
		configPath    string
		oldPrefix     string
		oldShort      string
		newPrefix     string
		newShort      string
		newCommand    string
		expectedError string
		expectedCmds  map[string]map[string]string
	}{
		{
			name: "Update successfully with new prefix, short, and command",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{
					"oldPrefix": {"oldShort": "original command"},
				},
			},
			oldPrefix:  "oldPrefix",
			oldShort:   "oldShort",
			newPrefix:  "newPrefix",
			newShort:   "newShort",
			newCommand: "updated command",
			expectedCmds: map[string]map[string]string{
				"newPrefix": {"newShort": "updated command"},
			},
		},
		{
			name: "Update without changes (no updates specified)",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{
					"oldPrefix": {"oldShort": "original command"},
				},
			},
			oldPrefix:  "oldPrefix",
			oldShort:   "oldShort",
			newPrefix:  "",
			newShort:   "",
			newCommand: "",
			expectedCmds: map[string]map[string]string{
				"oldPrefix": {"oldShort": "original command"},
			},
		},
		{
			name: "Fail update with nonexistent oldPrefix",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{},
			},
			oldPrefix:     "nonexistentPrefix",
			oldShort:      "oldShort",
			newPrefix:     "newPrefix",
			newShort:      "newShort",
			newCommand:    "new command",
			expectedError: "prefix not found: nonexistentPrefix",
		},
		{
			name: "Fail update with nonexistent oldShort",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{
					"oldPrefix": {},
				},
			},
			oldPrefix:     "oldPrefix",
			oldShort:      "nonexistentShort",
			newPrefix:     "newPrefix",
			newShort:      "newShort",
			newCommand:    "new command",
			expectedError: "command not found: oldPrefix nonexistentShort",
		},
		{
			name: "Update with only newPrefix",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{
					"oldPrefix": {"oldShort": "original command"},
				},
			},
			oldPrefix:  "oldPrefix",
			oldShort:   "oldShort",
			newPrefix:  "newPrefix",
			newShort:   "",
			newCommand: "",
			expectedCmds: map[string]map[string]string{
				"newPrefix": {"oldShort": "original command"},
			},
		},
		{
			name: "Update with only newShort",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{
					"oldPrefix": {"oldShort": "original command"},
				},
			},
			oldPrefix:  "oldPrefix",
			oldShort:   "oldShort",
			newPrefix:  "",
			newShort:   "newShort",
			newCommand: "",
			expectedCmds: map[string]map[string]string{
				"oldPrefix": {"newShort": "original command"},
			},
		},
		{
			name: "Update with only newCommand",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{
					"oldPrefix": {"oldShort": "original command"},
				},
			},
			oldPrefix:  "oldPrefix",
			oldShort:   "oldShort",
			newPrefix:  "",
			newShort:   "",
			newCommand: "updated command",
			expectedCmds: map[string]map[string]string{
				"oldPrefix": {"oldShort": "updated command"},
			},
		},
		{
			name: "Fail update with newShort already existing",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{
					"oldPrefix": {
						"oldShort": "original command",
						"newShort": "existing command",
					},
				},
			},
			oldPrefix:     "oldPrefix",
			oldShort:      "oldShort",
			newPrefix:     "",
			newShort:      "newShort",
			newCommand:    "",
			expectedError: "short already exists: newShort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "test_config_*.toml")
			if err != nil {
				t.Fatalf("failed to create temporary file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			if tt.initialConfig != nil {
				if err := config.SaveConfig(tt.initialConfig, tempFile.Name()); err != nil {
					t.Fatalf("failed to save initial config: %v", err)
				}
			}

			err = handler.UpdateCommand(tempFile.Name(), tt.oldPrefix, tt.oldShort, tt.newPrefix, tt.newShort, tt.newCommand)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("unexpected error: got %v, want %v", err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			loadedConfig, err := config.LoadConfig(tempFile.Name())
			if err != nil {
				t.Fatalf("failed to load config after execution: %v", err)
			}
			if !equalCommands(loadedConfig.Commands, tt.expectedCmds) {
				t.Errorf("unexpected commands in config: got %v, want %v", loadedConfig.Commands, tt.expectedCmds)
			}
		})
	}
}
