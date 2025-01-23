package handler

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
)

func ExecCommand(configPath, prefix, short string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	cmds, exists := cfg.Commands[prefix]
	if !exists {
		return fmt.Errorf("command not found: %s/%s", prefix, short)
	}

	command, ok := cmds[short]
	if !ok {
		return fmt.Errorf("command not found: %s/%s", prefix, short)
	}

	execCmd := exec.Command("sh", "-c", command)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin
	return execCmd.Run()
}
