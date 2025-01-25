package config

func (c *Config) GetRegisteredPrefixes() []string {
	prefixes := make([]string, 0, len(c.Commands))
	for p := range c.Commands {
		prefixes = append(prefixes, p)
	}
	return prefixes
}

func (c *Config) GetRegisteredShortcutsByPrefix(prefix string) []string {
	cmds, exists := c.Commands[prefix]
	if !exists {
		return nil
	}

	shorts := make([]string, 0, len(cmds))
	for s := range cmds {
		shorts = append(shorts, s)
	}
	return shorts
}
