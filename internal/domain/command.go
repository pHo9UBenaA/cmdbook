package domain

type CommandEntry struct {
	Prefix  string
	Short   string
	Command string
}

func GroupCommands(commands map[string]map[string]string) map[string][]CommandEntry {
	grouped := make(map[string][]CommandEntry)
	for prefix, cmds := range commands {
		var entries []CommandEntry
		for short, cmd := range cmds {
			entries = append(entries, CommandEntry{prefix, short, cmd})
		}
		grouped[prefix] = entries
	}
	return grouped
}

func PrepareInteractiveEntries(grouped map[string][]CommandEntry) []CommandEntry {
	var entries []CommandEntry
	for prefix, cmds := range grouped {
		entries = append(entries, CommandEntry{Prefix: prefix})
		entries = append(entries, cmds...)
	}
	return entries
}
