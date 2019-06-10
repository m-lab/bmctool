package cmd

import (
	"context"
	"fmt"

	"github.com/m-lab/go/rtx"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <host> [flags]",
	Short: "Get credentials for a BMC",
	Long: `This command gets a Credentials entity for a given BMC from Google
Cloud Datastore.

The GCP project to use can be specified by providing the --project flag.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		printCredentials(host)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

// printCredentials retrieves credentials for a given hostname and prints them
// in JSON format.
func printCredentials(host string) {
	provider := credsNewProvider(projectID, namespace)
	creds, err := provider.FindCredentials(context.Background(), makeBMCHostname(host))
	rtx.Must(err, "Cannot fetch credentials")

	fmt.Print(creds)
}
