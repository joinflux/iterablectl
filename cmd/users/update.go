package users

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joinflux/iterablectl/pkg/iterable"
	"github.com/spf13/cobra"
)

// UpdateCmd represents the update command for users
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a user in Iterable",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		client := iterable.NewClient(apiKey)

		email, _ := cmd.Flags().GetString("email")
		userId, _ := cmd.Flags().GetString("user-id")
		mergeNestedObjects, _ := cmd.Flags().GetBool("merge-nested-objects")

		// Either email or userId must be provided
		if email == "" && userId == "" {
			return fmt.Errorf("either --email or --user-id must be provided")
		}

		// Create user object
		user := iterable.User{
			Email:              email,
			UserID:             userId,
			DataFields:         make(map[string]interface{}),
			MergeNestedObjects: mergeNestedObjects,
		}

		// Handle data fields
		dataFieldsStr, _ := cmd.Flags().GetStringArray("data-field")
		for _, field := range dataFieldsStr {
			parts := strings.SplitN(field, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid data field format: %s, expected format is key=value", field)
			}

			key, value := parts[0], parts[1]
			user.DataFields[key] = value
		}

		// Handle data fields from file
		dataFile, _ := cmd.Flags().GetString("data-file")
		if dataFile != "" {
			fileData, err := os.ReadFile(dataFile)
			if err != nil {
				return fmt.Errorf("failed to read data file: %v", err)
			}

			var dataFields map[string]interface{}
			if err := json.Unmarshal(fileData, &dataFields); err != nil {
				return fmt.Errorf("failed to parse data file as JSON: %v", err)
			}

			// Merge file data with any command-line data fields
			for k, v := range dataFields {
				user.DataFields[k] = v
			}
		}

		// Update the user
		if err := client.UpdateUser(user); err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}

		fmt.Println("User updated successfully")
		return nil
	},
}

func init() {
	UpdateCmd.Flags().String("email", "", "User email address")
	UpdateCmd.Flags().String("user-id", "", "User ID")
	UpdateCmd.Flags().StringArray("data-field", []string{}, "Data field in key=value format (can be used multiple times)")
	UpdateCmd.Flags().String("data-file", "", "JSON file containing data fields")
	UpdateCmd.Flags().Bool("merge-nested-objects", true, "Whether to merge nested objects")
}
