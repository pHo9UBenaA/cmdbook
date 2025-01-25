package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{Commands: make(map[string]map[string]string)}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
