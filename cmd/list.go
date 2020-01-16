package cmd

import (
	"context"
	"fmt"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/creds"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List all the available BMCs.",
	Long: `This command lists all the BMCs stored on Google Cloud Datastore.

The GCP project to use can be specified by providing the --project flag`,
	Run: func(cmd *cobra.Command, args []string) {
		listBMCs()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listBMCs() {
	// If no project ID has been specified, don't do anything.
	if projectID == "" {
		osExit(1)
	}

	provider, err := credsNewProvider(&creds.DatastoreConnector{}, projectID, namespace)
	rtx.Must(err, "Cannot connect to Datastore")

	creds, err := provider.ListCredentials(context.Background())
	rtx.Must(err, "Cannot read credentials list")

	for _, c := range creds {
		fmt.Println(c.Hostname)
	}
}
