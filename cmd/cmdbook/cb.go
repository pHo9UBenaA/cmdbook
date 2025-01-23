package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

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

	rootCmd.AddCommand(
		addCmd(),
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

	cmd := &cobra.Command{
		Use:   "add <command>",
		Short: "Add a new command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.AddCommand(configPath, prefix, short, args[0]); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
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
			if err := handler.ExecCommand(configPath, args[0], args[1]); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}
}

func removeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <prefix> <short-cmd>",
		Short: "Remove a command",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.RemoveCommand(configPath, args[0], args[1]); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}
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
