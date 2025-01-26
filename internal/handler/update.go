package handler

import (
	"fmt"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
)

func UpdateCommand(configPath string, oldPrefix, oldShort, newPrefix, newShort, newCommand string) error {
	if newPrefix == "" && newShort == "" && newCommand == "" {
		fmt.Println("No updates specified. Skipping command update.")
		return nil
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cmds, exists := cfg.Commands[oldPrefix]
	if !exists {
		return fmt.Errorf("prefix not found: %s", oldPrefix)
	}

	originalCmd, ok := cmds[oldShort]
	if !ok {
		return fmt.Errorf("command not found: %s %s", oldPrefix, oldShort)
	}

	if newPrefix == "" {
		newPrefix = oldPrefix
	} else {
		err := updatePrefix(cfg, oldPrefix, oldShort, newPrefix, originalCmd)

		if err != nil {
			return err
		}
	}

	if newShort == "" {
		newShort = oldShort
	} else {
		err := updateShort(cfg, newPrefix, oldShort, newShort, originalCmd)

		if err != nil {
			return err
		}
	}

	if err := updateCommand(cfg, newPrefix, newShort, newCommand); err != nil {
		return err
	}

	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Updated: %s %s -> %s %s\n", oldPrefix, oldShort, newPrefix, newShort)
	return nil
}

func updatePrefix(cfg *config.Config, oldPrefix, oldShort, newPrefix, originalCmd string) error {
	cmds := cfg.Commands[oldPrefix]
	delete(cmds, oldShort)

	if len(cmds) == 0 {
		delete(cfg.Commands, oldPrefix)
	}

	if cfg.Commands[newPrefix] == nil {
		cfg.Commands[newPrefix] = make(map[string]string)
	}

	cfg.Commands[newPrefix][oldShort] = originalCmd

	return nil
}

func updateShort(cfg *config.Config, prefix, oldShort, newShort, originalCmd string) error {
	cmds := cfg.Commands[prefix]
	if _, ok := cmds[newShort]; ok {
		return fmt.Errorf("short already exists: %s", newShort)
	}
	delete(cmds, oldShort)

	cmds[newShort] = originalCmd

	return nil
}

func updateCommand(cfg *config.Config, prefix, short, newCommand string) error {
	if newCommand == "" {
		return nil
	}

	cmds := cfg.Commands[prefix]
	cmds[short] = newCommand

	return nil
}
