package cmd

import (
	"context"

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
	if projectID == "" {
		projectID = getProjectID(bmcHost)
	}

	log.Infof("Deleting credentials for host %v", bmcHost)
	provider := credsNewProvider(projectID, namespace)

	err := provider.DeleteCredentials(context.Background(), bmcHost)
	if err != nil {
		log.Errorf("Cannot delete credentials for %s: %v", bmcHost, err)
		osExit(1)
	}

	log.Info("Credentials successfully deleted.")
}
