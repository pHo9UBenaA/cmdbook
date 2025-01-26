package handler_test

import (
	"os"
	"testing"

	"github.com/pHo9UBenaA/cmdbook/internal/handler"
)

// Helper function to create a temporary configuration file
func createTempConfig(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "testconfig_*.toml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// Helper function to clean up temporary files
func cleanupTempFile(path string) {
	_ = os.Remove(path)
}

func TestExecCommand(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		configPath    string
		prefix        string
		short         string
		expectedError string
	}{
		{
			name: "successfully executes command",
			configContent: `
[commands.build]
docker = "echo building docker image"
`,
			prefix:        "build",
			short:         "docker",
			expectedError: "",
		},
		{
			name: "prefix not found in config",
			configContent: `
[commands.build]
docker = "echo building docker image"
`,
			prefix:        "deploy",
			short:         "k8s",
			expectedError: "command not found: deploy k8s",
		},
		{
			name: "short command not found in prefix",
			configContent: `
[commands.build]
docker = "echo building docker image"
`,
			prefix:        "build",
			short:         "java",
			expectedError: "command not found: build java",
		},
		{
			name:          "config file not found",
			configContent: "",
			prefix:        "build",
			short:         "docker",
			expectedError: "command not found: build docker",
		},
		{
			name: "exec command fails",
			configContent: `
[commands.build]
fail = "exit 1"
`,
			prefix:        "build",
			short:         "fail",
			expectedError: "exit status 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var configPath string
			var err error

			// Create a temporary config file if content is provided
			if tt.configContent != "" {
				configPath, err = createTempConfig(tt.configContent)
				if err != nil {
					t.Fatalf("failed to create temp config file: %v", err)
				}
				defer cleanupTempFile(configPath)
			} else {
				// Use a non-existent file path for the "file not found" case
				configPath = "/nonexistent/config.toml"
			}

			// Execute the function
			err = handler.ExecCommand(configPath, tt.prefix, tt.short)

			// Validate the error
			if (err != nil && err.Error() != tt.expectedError) || (err == nil && tt.expectedError != "") {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedError)
			}
		})
	}
}
