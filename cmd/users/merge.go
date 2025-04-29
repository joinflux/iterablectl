package users

import (
	"fmt"

	"github.com/joinflux/iterablectl/pkg/iterable"
	"github.com/spf13/cobra"
)

// MergeCmd represents the merge command for users
var MergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge two users",
	Example: `iterablectl users merge --from-email <source email> --to-email <destination email>
iterablectl users merge --from-user-id <source user id> --to-email <destination email>
iterablectl users merge --from-user-id <source user id> --to-user-id <destination user id>
iterablectl users merge --from-email <source email> --to-user-id <destination user id>
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		client := iterable.NewClient(apiKey)

		// Get flag values
		fromEmail, _ := cmd.Flags().GetString("from-email")
		fromUserID, _ := cmd.Flags().GetString("from-user-id")
		toEmail, _ := cmd.Flags().GetString("to-email")
		toUserID, _ := cmd.Flags().GetString("to-user-id")

		// Validate that exactly one source identifier is provided
		if (fromEmail == "" && fromUserID == "") || (fromEmail != "" && fromUserID != "") {
			return fmt.Errorf("exactly one of --source-email or --source-user-id must be specified")
		}
		if (toEmail == "" && toUserID == "") || (toEmail != "" && toUserID != "") {
			return fmt.Errorf("exactly one of --to-email or --to-user-id must be specified")
		}

		var response *iterable.APIError
		var err error

		response, err = client.MergeUsers(iterable.MergeUsersOpts{
			SrcEmail: fromEmail,
			DstEmail: toEmail,

			SrcID: fromUserID,
			DstID: toUserID,
		})
		if err != nil {
			return fmt.Errorf("error merging users: %w", err)
		}

		if response.Code != "Success" {
			return fmt.Errorf("error merging users: %v", response)
		}

		return nil
	},
}

func init() {
	MergeCmd.Flags().String("from-email", "", "Email address of the source user to merge")
	MergeCmd.Flags().String("from-user-id", "", "User ID of the source user to merge")
	MergeCmd.Flags().String("to-email", "", "Email address of the destination profile to merge into")
	MergeCmd.Flags().String("to-user-id", "", "User ID of the destination profile to merge into")
}
