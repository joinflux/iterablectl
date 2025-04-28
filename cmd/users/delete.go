package users

import (
	"fmt"

	"github.com/joinflux/iterablectl/pkg/iterable"
	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command for users
var DeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a user from Iterable",
	Args:    cobra.ExactArgs(1),
	Example: "iterablectl users delete user@example.com",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		client := iterable.NewClient(apiKey)

		email := args[0]
		if email == "" {
			return fmt.Errorf("email is required")
		}

		byUserID, _ := cmd.Flags().GetBool("by-userid")

		var err error
		if byUserID {
			err = client.DeleteUserByID(email)
		} else {
			err = client.DeleteUser(email)
		}

		if err != nil {
			return fmt.Errorf("error deleting user: %v", err)
		}

		identifier := "email"
		if byUserID {
			identifier = "userID"
		}
		fmt.Printf("User with %s '%s' successfully deleted\n", identifier, email)
		return nil
	},
}

func init() {
	DeleteCmd.Flags().Bool("by-userid", false, "Delete user by user ID instead of email")
}
