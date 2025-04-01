package lists

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/joinflux/iterablectl/pkg/iterable"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all lists from Iterable",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		client := iterable.NewClient(apiKey)

		lists, err := client.GetLists()
		if err != nil {
			return fmt.Errorf("error getting user: %v", err)
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "json" {
			jsonOutput, err := json.MarshalIndent(lists, "", "  ")
			if err != nil {
				return fmt.Errorf("error formatting JSON: %v", err)
			}
			fmt.Println(string(jsonOutput))
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, list := range *lists {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", list.ID, list.Name, list.Description, time.Unix(list.CreatedAt/1000, 0).Local().Format(time.DateOnly), list.ListType)
		}
		w.Flush()

		return nil
	},
}
