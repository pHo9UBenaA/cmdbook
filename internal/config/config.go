package config

type Config struct {
	Commands map[string]map[string]string `toml:"commands"`
}
