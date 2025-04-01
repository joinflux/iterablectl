package lists

import (
	"github.com/spf13/cobra"
)

// Cmd represents the lists command
var Cmd = &cobra.Command{
	Use:   "lists",
	Short: "Manage Iterable lists",
}

func init() {
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(UsersCmd)
}
