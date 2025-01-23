package handler

import (
	"fmt"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
)

func RemoveCommand(configPath, prefix, shortCmd string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cmds, exists := cfg.Commands[prefix]
	if !exists {
		return fmt.Errorf("prefix does not exist: %s", prefix)
	}

	if _, ok := cmds[shortCmd]; !ok {
		return fmt.Errorf("command not found: %s/%s", prefix, shortCmd)
	}

	delete(cmds, shortCmd)
	if len(cmds) == 0 {
		delete(cfg.Commands, prefix)
	}

	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Removed: %s/%s\n", prefix, shortCmd)
	return nil
}
