package campaigns

import (
	"github.com/spf13/cobra"
)

// Cmd represents the campaigns command
var Cmd = &cobra.Command{
	Use:   "campaigns",
	Short: "Manage Iterable campaigns",
}

func init() {
	Cmd.AddCommand(GetCmd)
}
