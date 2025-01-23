package handler

import (
	"fmt"
	"strings"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
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
		for i := 0; ; i++ {
			candidate := fmt.Sprintf("cmd%d", i)
			if _, exists := existing[candidate]; !exists {
				short = candidate
				break
			}
		}
	}

	if cfg.Commands[prefix] == nil {
		cfg.Commands[prefix] = make(map[string]string)
	}

	cfg.Commands[prefix][short] = command
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return err
	}

	fmt.Printf("Added: %s/%s -> %s\n", prefix, short, command)
	return nil
}
