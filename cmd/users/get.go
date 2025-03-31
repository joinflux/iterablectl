package users

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/joinflux/iterablectl/pkg/iterable"
	"github.com/joinflux/iterablectl/pkg/utils"
	"github.com/spf13/cobra"
)

// GetCmd represents the get command for users
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a user from Iterable by email",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		client := iterable.NewClient(apiKey)

		email := args[0]
		if email == "" {
			return fmt.Errorf("email is required")
		}

		user, err := client.GetUser(email)
		if err != nil {
			return fmt.Errorf("error getting user: %v", err)
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "json" {
			jsonOutput, err := json.MarshalIndent(user, "", "  ")
			if err != nil {
				return fmt.Errorf("error formatting JSON: %v", err)
			}
			fmt.Println(string(jsonOutput))
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		printNestedFields("", user.DataFields, w)
		w.Flush()

		return nil
	},
}

// printNestedFields recursively prints nested data fields
func printNestedFields(prefix string, data map[string]any, w *tabwriter.Writer) {
	// Sort keys for consistent output
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		key := prefix
		if prefix != "" {
			key = prefix + "." + k
		} else {
			key = k
		}

		switch value := v.(type) {
		case map[string]any:
			// Recursively print nested objects
			printNestedFields(key, value, w)
		default:
			// Print other values
			fmt.Fprintf(w, "%s\t%v\n", key, utils.FormatValue(v))
		}
	}
}

func init() {
	GetCmd.Flags().String("format", "table", "Output format: json or table (default)")
}
