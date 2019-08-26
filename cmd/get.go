package cmd

import (
	"context"
	"fmt"

	"github.com/m-lab/reboot-service/creds"

	"github.com/m-lab/go/rtx"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <host> [flags]",
	Short: "Get credentials for a BMC",
	Long: `This command gets a Credentials entity for a given BMC from Google
Cloud Datastore.

The GCP project to use can be specified by providing the --project flag,
otherwise it will be inferred by the node name.`,
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
	bmcHost := makeBMCHostname(host)
	if projectID == "" {
		projectID = getProjectID(bmcHost)
	}

	provider, err := credsNewProvider(&creds.DatastoreConnector{}, projectID, namespace)
	rtx.Must(err, "Cannot connect to Datastore")
	defer provider.Close()

	creds, err := provider.FindCredentials(context.Background(), bmcHost)
	rtx.Must(err, "Cannot fetch credentials")

	fmt.Print(creds)
}
