package config_test

import (
	"sort"
	"testing"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
)

func TestConfig_GetRegisteredPrefixes(t *testing.T) {
	tests := []struct {
		name     string
		commands map[string]map[string]string
		want     []string
	}{
		{
			name:     "empty commands",
			commands: map[string]map[string]string{},
			want:     []string{},
		},
		{
			name:     "single prefix",
			commands: map[string]map[string]string{"p1": {}},
			want:     []string{"p1"},
		},
		{
			name: "multiple prefixes",
			commands: map[string]map[string]string{
				"p3": {},
				"p1": {},
				"p2": {},
			},
			want: []string{"p1", "p2", "p3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Commands: tt.commands,
			}
			got := cfg.GetRegisteredPrefixes()
			sort.Strings(got)
			sort.Strings(tt.want)

			if len(got) != len(tt.want) {
				t.Fatalf("got %d prefixes (%v), want %d (%v)", len(got), got, len(tt.want), tt.want)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("prefix mismatch at index %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestConfig_GetRegisteredShortcutsByPrefix(t *testing.T) {
	tests := []struct {
		name     string
		commands map[string]map[string]string
		prefix   string
		want     []string
	}{
		{
			name: "prefix not found",
			commands: map[string]map[string]string{
				"other": {"key": "value"},
			},
			prefix: "test",
			want:   nil,
		},
		{
			name: "prefix exists with no shortcuts",
			commands: map[string]map[string]string{
				"test": {},
			},
			prefix: "test",
			want:   []string{},
		},
		{
			name: "prefix exists with nil shortcuts map",
			commands: map[string]map[string]string{
				"test": nil,
			},
			prefix: "test",
			want:   []string{},
		},
		{
			name: "prefix exists with one shortcut",
			commands: map[string]map[string]string{
				"test": {"s1": "cmd"},
			},
			prefix: "test",
			want:   []string{"s1"},
		},
		{
			name: "prefix exists with multiple shortcuts",
			commands: map[string]map[string]string{
				"test": {"s2": "cmd", "s1": "cmd", "s3": "cmd"},
			},
			prefix: "test",
			want:   []string{"s1", "s2", "s3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Commands: tt.commands,
			}
			got := cfg.GetRegisteredShortcutsByPrefix(tt.prefix)

			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}

			sort.Strings(got)
			sort.Strings(tt.want)

			if len(got) != len(tt.want) {
				t.Fatalf("got %d shortcuts (%v), want %d (%v)", len(got), got, len(tt.want), tt.want)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("shortcut mismatch at index %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
