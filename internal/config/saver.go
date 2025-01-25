package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func SaveConfig(cfg *Config, path string) error {
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
