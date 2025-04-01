package users

import (
	"fmt"

	"github.com/joinflux/iterablectl/pkg/iterable"
	"github.com/spf13/cobra"
)

// MergeCmd represents the merge command for users
var MergeCmd = &cobra.Command{
	Use:     "merge",
	Short:   "Merge two users",
	Args:    cobra.ExactArgs(2),
	Example: "iterablectl users merge <source_email> <destination_email>",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		client := iterable.NewClient(apiKey)

		sourceEmail := args[0]
		if sourceEmail == "" {
			return fmt.Errorf("source email is required")
		}
		destinationEmail := args[1]
		if destinationEmail == "" {
			return fmt.Errorf("destination email is required")
		}

		response, err := client.MergeUsers(sourceEmail, destinationEmail)
		if err != nil {
			return fmt.Errorf("error merging users: %v", err)
		}

		if response.Code != "Success" {
			return fmt.Errorf("error merging users: %v", response)
		}

		return nil
	},
}
