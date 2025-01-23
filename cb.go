package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type CommandEntry struct {
	Prefix  string
	Short   string
	Command string
}

type Config struct {
	Commands map[string]map[string]string `toml:"commands"`
}

var configPath string

func main() {
	home, _ := os.UserHomeDir()
	configPath = filepath.Join(home, ".cmdbook.toml")

	rootCmd := &cobra.Command{
		Use:   "cb",
		Short: "Command Book - Manage your frequently used commands",
	}

	rootCmd.AddCommand(
		addCmd(),
		execCmd(),
		removeCmd(),
		listCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func loadConfig() (*Config, error) {
	config := &Config{Commands: make(map[string]map[string]string)}

	if data, err := os.ReadFile(configPath); err == nil {
		if err := toml.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func saveConfig(config *Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func addCmd() *cobra.Command {
	var short, prefix string

	cmd := &cobra.Command{
		Use:   "add <command>",
		Short: "Add a new command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, _ := loadConfig()
			command := args[0]

			if prefix == "" {
				if parts := strings.SplitN(command, " ", 2); len(parts) > 0 {
					prefix = parts[0]
				}
			}

			if short == "" {
				existing := config.Commands[prefix]
				for i := 0; ; i++ {
					candidate := fmt.Sprintf("cmd%d", i)
					if _, exists := existing[candidate]; !exists {
						short = candidate
						break
					}
				}
			}

			if config.Commands[prefix] == nil {
				config.Commands[prefix] = make(map[string]string)
			}

			config.Commands[prefix][short] = command
			saveConfig(config)

			fmt.Printf("Added: %s/%s -> %s\n", prefix, short, command)
		},
	}

	cmd.Flags().StringVarP(&short, "short", "s", "", "Short command name")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Command prefix")

	return cmd
}

func execCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exec <prefix> <short-cmd>",
		Short: "Execute a command",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			config, _ := loadConfig()
			prefix, short := args[0], args[1]

			if cmds, exists := config.Commands[prefix]; exists {
				if command, ok := cmds[short]; ok {
					execCmd := exec.Command("sh", "-c", command)
					execCmd.Stdout = os.Stdout
					execCmd.Stderr = os.Stderr
					execCmd.Stdin = os.Stdin
					execCmd.Run()
					return
				}
			}
			fmt.Printf("Command not found: %s/%s\n", prefix, short)
		},
	}
}

func removeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <prefix> <short-cmd>",
		Short: "Remove a command",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			config, _ := loadConfig()
			prefix, short := args[0], args[1]

			if cmds, exists := config.Commands[prefix]; exists {
				if _, ok := cmds[short]; ok {
					delete(cmds, short)
					if len(cmds) == 0 {
						delete(config.Commands, prefix)
					}
					saveConfig(config)
					fmt.Printf("Removed: %s/%s\n", prefix, short)
					return
				}
			}
			fmt.Printf("Command not found: %s/%s\n", prefix, short)
		},
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List commands with interactive viewer",
		Run: func(cmd *cobra.Command, args []string) {
			config, _ := loadConfig()
			var entries []CommandEntry

			for prefix, cmds := range config.Commands {
				for short, command := range cmds {
					entries = append(entries, CommandEntry{prefix, short, command})
				}
			}

			sort.Slice(entries, func(i, j int) bool {
				if entries[i].Prefix == entries[j].Prefix {
					return entries[i].Short < entries[j].Short
				}
				return entries[i].Prefix < entries[j].Prefix
			})

			if len(entries) == 0 {
				fmt.Println("No commands stored")
				return
			}

			height, _, _ := term.GetSize(int(os.Stdout.Fd()))
			pageSize := height - 2
			if pageSize < 1 {
				pageSize = 10
			}

			keyboard.Open()
			defer keyboard.Close()

			offset := 0
			for {
				fmt.Print("\033[2J\033[H") // Clear screen
				for i := offset; i < offset+pageSize && i < len(entries); i++ {
					fmt.Printf("%s/%s: %s\n", entries[i].Prefix, entries[i].Short, entries[i].Command)
				}

				fmt.Printf("\nCommands %d-%d of %d (▲/▼ scroll, q quit)", offset+1, min(offset+pageSize, len(entries)), len(entries))

				char, key, _ := keyboard.GetKey()
				switch {
				case key == keyboard.KeyArrowUp && offset > 0:
					offset--
				case key == keyboard.KeyArrowDown && offset < len(entries)-pageSize:
					offset++
				case char == 'q' || key == keyboard.KeyEsc:
					return
				}
			}
		},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
