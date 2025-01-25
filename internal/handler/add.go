package handler

import (
	"fmt"
	"strings"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
	"github.com/pHo9UBenaA/cmdbook/internal/constant"
)

func AddCommand(configPath, prefix, short, command string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	if prefix == "" {
		prefix = strings.SplitN(command, " ", 2)[0]
	}

	if short == "" {
		existing := cfg.Commands[prefix]
		nextIndex := len(existing)
		short = fmt.Sprintf("cmd%d", nextIndex)
	}

	if len(short) > constant.MaxShortLen {
		return fmt.Errorf("short name '%s' exceeds maximum length of 20 characters", short)
	}

	if cfg.Commands[prefix] == nil {
		cfg.Commands[prefix] = make(map[string]string)
	}

	cfg.Commands[prefix][short] = command
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return err
	}

	fmt.Printf("Added: %s %s -> %s ", prefix, short, command)
	return nil
}
