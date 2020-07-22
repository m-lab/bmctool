package cmd

import (
	"context"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/creds"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete <hostname>",
		Short: "Deletes the Credentials entity for the specified hostname",
		Long:  "This command deletes the Credentials entity on GCD.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bmcHost = args[0]
			deleteCredentials()
		},
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

// deleteCredentials updates a Credentials entity on Google Cloud Datastore.
func deleteCredentials() {
	bmcHost = makeBMCHostname(bmcHost)

	log.Infof("Deleting credentials for host %v", bmcHost)
	provider, err := credsNewProvider(&creds.DatastoreConnector{}, projectID, namespace)
	rtx.Must(err, "Cannot connect to Datastore")
	provider.Close()

	err = provider.DeleteCredentials(context.Background(), bmcHost)
	// Note: Deleting a key from Datastore does not return a NoSuchEntity error
	// if the specified key does not exist, thus the error will be nil unless
	// something else goes wrong during the deletion.
	// See: https://github.com/googleapis/google-cloud-go/issues/501
	if err != nil {
		log.Errorf("Cannot delete credentials for %s: %v", bmcHost, err)
		osExit(1)
	}

	log.Info("Credentials successfully deleted.")
}
