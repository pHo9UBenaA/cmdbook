package domain

import (
	"sort"
	"testing"
)

func TestGroupCommands(t *testing.T) {
	tests := []struct {
		name     string
		commands map[string]map[string]string
		expected map[string][]CommandEntry
	}{
		{
			name:     "empty commands",
			commands: map[string]map[string]string{},
			expected: map[string][]CommandEntry{},
		},
		{
			name: "single prefix with one shortcut",
			commands: map[string]map[string]string{
				"p1": {"s1": "cmd1"},
			},
			expected: map[string][]CommandEntry{
				"p1": {{Prefix: "p1", Short: "s1", Command: "cmd1"}},
			},
		},
		{
			name: "single prefix with multiple shortcuts",
			commands: map[string]map[string]string{
				"p1": {"s2": "cmd2", "s1": "cmd1"},
			},
			expected: map[string][]CommandEntry{
				"p1": {
					{Prefix: "p1", Short: "s1", Command: "cmd1"},
					{Prefix: "p1", Short: "s2", Command: "cmd2"},
				},
			},
		},
		{
			name: "multiple prefixes",
			commands: map[string]map[string]string{
				"p1": {"s1": "cmd1"},
				"p2": {"s3": "cmd3"},
			},
			expected: map[string][]CommandEntry{
				"p1": {{Prefix: "p1", Short: "s1", Command: "cmd1"}},
				"p2": {{Prefix: "p2", Short: "s3", Command: "cmd3"}},
			},
		},
		{
			name: "prefix with empty shortcuts map",
			commands: map[string]map[string]string{
				"p1": {},
			},
			expected: map[string][]CommandEntry{
				"p1": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupCommands(tt.commands)

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d groups, got %d", len(tt.expected), len(got))
			}

			for expPrefix, expEntries := range tt.expected {
				actEntries, exists := got[expPrefix]
				if !exists {
					t.Errorf("expected prefix %q not found", expPrefix)
					continue
				}

				sort.Slice(actEntries, func(i, j int) bool {
					return actEntries[i].Short < actEntries[j].Short
				})
				sort.Slice(expEntries, func(i, j int) bool {
					return expEntries[i].Short < expEntries[j].Short
				})

				if len(actEntries) != len(expEntries) {
					t.Errorf("prefix %q: expected %d entries, got %d", expPrefix, len(expEntries), len(actEntries))
					continue
				}

				for i := range actEntries {
					if actEntries[i] != expEntries[i] {
						t.Errorf("prefix %q entry %d mismatch: got %v, want %v",
							expPrefix, i, actEntries[i], expEntries[i])
					}
				}
			}
		})
	}
}

func TestPrepareInteractiveEntries(t *testing.T) {
	tests := []struct {
		name           string
		grouped        map[string][]CommandEntry
		expectedOrder  []CommandEntry // Headers in expected order
		expectedGroups map[string][]CommandEntry
	}{
		{
			name:           "empty grouped",
			grouped:        map[string][]CommandEntry{},
			expectedOrder:  []CommandEntry{},
			expectedGroups: map[string][]CommandEntry{},
		},
		{
			name: "single prefix with no commands",
			grouped: map[string][]CommandEntry{
				"p1": {},
			},
			expectedOrder: []CommandEntry{{Prefix: "p1"}},
			expectedGroups: map[string][]CommandEntry{
				"p1": {{Prefix: "p1"}},
			},
		},
		{
			name: "single prefix with nil commands slice",
			grouped: map[string][]CommandEntry{
				"p1": nil,
			},
			expectedOrder: []CommandEntry{{Prefix: "p1"}},
			expectedGroups: map[string][]CommandEntry{
				"p1": {{Prefix: "p1"}},
			},
		},
		{
			name: "single prefix with multiple commands",
			grouped: map[string][]CommandEntry{
				"p1": {
					{Prefix: "p1", Short: "s2", Command: "cmd2"},
					{Prefix: "p1", Short: "s1", Command: "cmd1"},
				},
			},
			expectedOrder: []CommandEntry{
				{Prefix: "p1"},
				{Prefix: "p1", Short: "s1", Command: "cmd1"},
				{Prefix: "p1", Short: "s2", Command: "cmd2"},
			},
			expectedGroups: map[string][]CommandEntry{
				"p1": {
					{Prefix: "p1"},
					{Prefix: "p1", Short: "s1", Command: "cmd1"},
					{Prefix: "p1", Short: "s2", Command: "cmd2"},
				},
			},
		},
		{
			name: "multiple prefixes",
			grouped: map[string][]CommandEntry{
				"p1": {
					{Prefix: "p1", Short: "s1", Command: "cmd1"},
				},
				"p2": {
					{Prefix: "p2", Short: "s3", Command: "cmd3"},
				},
			},
			expectedGroups: map[string][]CommandEntry{
				"p1": {
					{Prefix: "p1"},
					{Prefix: "p1", Short: "s1", Command: "cmd1"},
				},
				"p2": {
					{Prefix: "p2"},
					{Prefix: "p2", Short: "s3", Command: "cmd3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrepareInteractiveEntries(tt.grouped)

			// Build actual groups from the output
			actualGroups := make(map[string][]CommandEntry)
			var currentPrefix string
			for _, entry := range got {
				if entry.Short == "" && entry.Command == "" {
					currentPrefix = entry.Prefix
					actualGroups[currentPrefix] = append(actualGroups[currentPrefix], entry)
				} else {
					actualGroups[currentPrefix] = append(actualGroups[currentPrefix], entry)
				}
			}

			// Verify expected groups
			for expPrefix, expGroup := range tt.expectedGroups {
				actGroup, exists := actualGroups[expPrefix]
				if !exists {
					t.Errorf("expected prefix %q not found in output", expPrefix)
					continue
				}

				if len(actGroup) == 0 {
					t.Errorf("prefix %q has no entries", expPrefix)
					continue
				}

				// Verify header
				header := actGroup[0]
				if header.Prefix != expPrefix || header.Short != "" || header.Command != "" {
					t.Errorf("prefix %q header mismatch: got %v", expPrefix, header)
				}

				// Verify commands (skip header)
				actCommands := actGroup[1:]
				expCommands := expGroup[1:]

				sort.Slice(actCommands, func(i, j int) bool {
					return actCommands[i].Short < actCommands[j].Short
				})
				sort.Slice(expCommands, func(i, j int) bool {
					return expCommands[i].Short < expCommands[j].Short
				})

				if len(actCommands) != len(expCommands) {
					t.Errorf("prefix %q: expected %d commands, got %d", expPrefix, len(expCommands), len(actCommands))
					continue
				}

				for i := range actCommands {
					if actCommands[i] != expCommands[i] {
						t.Errorf("prefix %q command %d mismatch: got %v, want %v",
							expPrefix, i, actCommands[i], expCommands[i])
					}
				}
			}

			// Verify no extra groups exist
			for actPrefix := range actualGroups {
				if _, exists := tt.expectedGroups[actPrefix]; !exists {
					t.Errorf("unexpected prefix %q in output", actPrefix)
				}
			}
		})
	}
}
