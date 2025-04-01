package campaigns

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/joinflux/iterablectl/pkg/iterable"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all campaigns from Iterable",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		client := iterable.NewClient(apiKey)

		campaigns, err := client.GetCampaigns()
		if err != nil {
			return fmt.Errorf("error getting campaigns: %v", err)
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "json" {
			jsonOutput, err := json.MarshalIndent(campaigns, "", "  ")
			if err != nil {
				return fmt.Errorf("error formatting JSON: %v", err)
			}
			fmt.Println(string(jsonOutput))
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, list := range *campaigns {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", list.ID, list.Name, list.CampaignState, list.MessageMedium)
		}
		w.Flush()

		return nil
	},
}
