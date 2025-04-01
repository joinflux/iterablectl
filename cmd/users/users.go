package users

import (
	"github.com/spf13/cobra"
)

// Cmd represents the user command
var Cmd = &cobra.Command{
	Use:   "users",
	Short: "Manage Iterable users",
	Long:  `Commands to list, update, and manage Iterable users.`,
}

func init() {
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(MergeCmd)
}
