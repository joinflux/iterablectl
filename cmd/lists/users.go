package lists

import (
	"fmt"

	"github.com/joinflux/iterablectl/pkg/iterable"
	"github.com/spf13/cobra"
)

var UsersCmd = &cobra.Command{
	Use:     "users",
	Short:   "Get users in a list",
	Args:    cobra.ExactArgs(1),
	Example: `iterablectl lists users <listId> [--ids true|false]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		client := iterable.NewClient(apiKey)

		listId := args[0]
		if listId == "" {
			return fmt.Errorf("listId is required")
		}

		preferUserId, _ := cmd.Flags().GetBool("ids")

		users, err := client.GetListUsers(listId, preferUserId)
		if err != nil {
			return fmt.Errorf("error getting users in list: %v", err)
		}
		fmt.Println(string(*users))

		return nil
	},
}

func init() {
	UsersCmd.Flags().BoolP("ids", "i", false, "Prefer userIds over email")
}
