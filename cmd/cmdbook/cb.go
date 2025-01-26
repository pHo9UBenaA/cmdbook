package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
	"github.com/pHo9UBenaA/cmdbook/internal/handler"
)

var configPath string

func main() {
	home, _ := os.UserHomeDir()
	configPath = filepath.Join(home, ".cmdbook.toml")

	rootCmd := &cobra.Command{
		Use:   "cb",
		Short: "Command Book - Manage your frequently used commands",
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.AddCommand(
		addCmd(),
		updateCmd(),
		execCmd(),
		removeCmd(),
		listCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addCmd() *cobra.Command {
	var short, prefix string

	const commandIndex = 0

	cmd := &cobra.Command{
		Use:   "add <command>",
		Short: "Add a new command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.AddCommand(configPath, prefix, short, args[commandIndex]); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&short, "short", "S", "", "Short command name")
	cmd.Flags().StringVarP(&prefix, "prefix", "P", "", "Command prefix")

	cmd.RegisterFlagCompletionFunc("prefix", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPrefixes(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func updateCmd() *cobra.Command {
	var newPrefix, newShort, newCommand string

	const (
		oldPrefixIndex   = 0
		oldShortCmdIndex = 1
		argsNum          = 2
	)

	cmd := &cobra.Command{
		Use:   "update <prefix> <short>",
		Short: "Update an existing command",
		Args:  cobra.ExactArgs(argsNum),
		Run: func(cmd *cobra.Command, args []string) {
			oldPrefix := args[oldPrefixIndex]
			oldShort := args[oldShortCmdIndex]

			if cmd.Flags().Changed("new-command") && newCommand == "" {
				fmt.Println("Error: new command cannot be empty")
				os.Exit(1)
			}

			if err := handler.UpdateCommand(configPath, oldPrefix, oldShort, newPrefix, newShort, newCommand); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&newPrefix, "new-prefix", "P", "", "New prefix for the command")
	cmd.Flags().StringVarP(&newShort, "new-short", "S", "", "New short name for the command")
	cmd.Flags().StringVarP(&newCommand, "new-command", "C", "", "New command content")

	cmd.RegisterFlagCompletionFunc("new-prefix", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPrefixes(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == oldPrefixIndex {
			return getPrefixes(), cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == oldShortCmdIndex {
			return getShorts(args[oldPrefixIndex]), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

func execCmd() *cobra.Command {
	const (
		prefixIndex   = 0
		shortCmdIndex = 1
		argsNum       = 2
	)

	cmd := &cobra.Command{
		Use:   "exec <prefix> <short-cmd>",
		Short: "Execute a command",
		Args:  cobra.ExactArgs(argsNum),
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.ExecCommand(configPath, args[prefixIndex], args[shortCmdIndex]); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == prefixIndex {
			return getPrefixes(), cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == shortCmdIndex {
			return getShorts(args[prefixIndex]), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

func removeCmd() *cobra.Command {
	const (
		prefixIndex   = 0
		shortCmdIndex = 1
		argsNum       = 2
	)

	cmd := &cobra.Command{
		Use:   "remove <prefix> <short-cmd>",
		Short: "Remove a command",
		Args:  cobra.ExactArgs(argsNum),
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.RemoveCommand(configPath, args[prefixIndex], args[shortCmdIndex]); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == prefixIndex {
			return getPrefixes(), cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == shortCmdIndex {
			return getShorts(args[prefixIndex]), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List commands with interactive viewer",
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.ListCommands(configPath); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}
}

func getPrefixes() []string {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil
	}

	return cfg.GetRegisteredPrefixes()
}

func getShorts(prefix string) []string {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil
	}

	return cfg.GetRegisteredShortcutsByPrefix(prefix)
}
