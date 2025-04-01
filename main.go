package main

import (
	"fmt"
	"os"

	"github.com/joinflux/iterablectl/cmd/lists"
	"github.com/joinflux/iterablectl/cmd/users"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "iterablectl",
	Short: "iterablectl - A command-line tool for Iterable API",
	Long: `iterablectl is a CLI tool that allows you to interface with the Iterable API.
You can list users, update user profiles, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommands are provided, show help
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("api-key", "k", os.Getenv("ITERABLE_API_KEY"), "Iterable API key (can also be set via ITERABLE_API_KEY environment variable)")

	// Only mark required if it's not set via environment variable
	if os.Getenv("ITERABLE_API_KEY") == "" {
		rootCmd.MarkPersistentFlagRequired("api-key")
	}

	// Add subcommands
	rootCmd.AddCommand(users.Cmd)
	rootCmd.AddCommand(lists.Cmd)
}

func main() {
	Execute()
}
