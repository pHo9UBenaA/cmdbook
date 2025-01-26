package handler_test

import (
	"os"
	"testing"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
	"github.com/pHo9UBenaA/cmdbook/internal/handler"
)

func TestAddCommand(t *testing.T) {
	tests := []struct {
		name          string
		initialConfig *config.Config
		configPath    string
		prefix        string
		short         string
		command       string
		expectedError error
		expectedCmds  map[string]map[string]string
	}{
		{
			name:          "Add command successfully with full input",
			initialConfig: &config.Config{Commands: map[string]map[string]string{}},
			prefix:        "testPrefix",
			short:         "testShort",
			command:       "echo Hello",
			expectedCmds: map[string]map[string]string{
				"testPrefix": {"testShort": "echo Hello"},
			},
		},
		{
			name:          "Add command with empty prefix, uses first word of command",
			initialConfig: &config.Config{Commands: map[string]map[string]string{}},
			prefix:        "",
			short:         "testShort",
			command:       "echo Hello",
			expectedCmds: map[string]map[string]string{
				"echo": {"testShort": "echo Hello"},
			},
		},
		{
			name: "Add command with empty short, auto-generate unique short",
			initialConfig: &config.Config{
				Commands: map[string]map[string]string{
					"testPrefix": {"cmd0": "some command"},
				},
			},
			prefix:  "testPrefix",
			short:   "",
			command: "echo Hello",
			expectedCmds: map[string]map[string]string{
				"testPrefix": {
					"cmd0": "some command",
					"cmd1": "echo Hello",
				},
			},
		},
		{
			name:          "Add command to non-existent config file",
			initialConfig: nil,
			prefix:        "newPrefix",
			short:         "newShort",
			command:       "new command",
			expectedCmds: map[string]map[string]string{
				"newPrefix": {"newShort": "new command"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file to act as the config file
			tempFile, err := os.CreateTemp("", "test_config_*.toml")
			if err != nil {
				t.Fatalf("failed to create temporary file: %v", err)
			}
			defer os.Remove(tempFile.Name()) // Clean up after the test

			// Write initial configuration to the file, if any
			if tt.initialConfig != nil {
				if err := config.SaveConfig(tt.initialConfig, tempFile.Name()); err != nil {
					t.Fatalf("failed to save initial config: %v", err)
				}
			}

			// Run the function under test
			err = handler.AddCommand(tempFile.Name(), tt.prefix, tt.short, tt.command)

			// Check for expected error
			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedError)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("unexpected error message: got %v, want %v", err.Error(), tt.expectedError.Error())
			}

			// Verify the config file contents
			loadedConfig, err := config.LoadConfig(tempFile.Name())
			if err != nil {
				t.Fatalf("failed to load config after execution: %v", err)
			}
			if len(loadedConfig.Commands) == 0 && tt.initialConfig == nil {
				loadedConfig.Commands = map[string]map[string]string{}
			}
			if !equalCommands(loadedConfig.Commands, tt.expectedCmds) {
				t.Errorf("unexpected commands in config: got %v, want %v", loadedConfig.Commands, tt.expectedCmds)
			}
		})
	}
}

// Helper function to compare two nested maps
func equalCommands(a, b map[string]map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, subMapA := range a {
		subMapB, ok := b[key]
		if !ok {
			return false
		}
		if len(subMapA) != len(subMapB) {
			return false
		}
		for subKey, valueA := range subMapA {
			valueB, ok := subMapB[subKey]
			if !ok || valueA != valueB {
				return false
			}
		}
	}
	return true
}
