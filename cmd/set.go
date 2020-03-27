package cmd

import (
	"context"
	"fmt"

	"github.com/m-lab/reboot-service/creds"

	"github.com/m-lab/go/rtx"
	"github.com/spf13/viper"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var (
	setCmd = &cobra.Command{
		Use:   "set <hostname>",
		Short: "Updates the Credentials entity for the specified hostname",
		Long: `This command updates the Credentials entity on GCD. If the flag
corresponding to a field is not specified, that field will not be updated.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bmcHost = args[0]
			setCredentials()
		},
	}
)

func init() {
	rootCmd.AddCommand(setCmd)

	viper.AutomaticEnv()

	setCmd.Flags().StringVar(&bmcUser, "user", viper.GetString("BMCUSER"),
		"BMC username")
	setCmd.Flags().StringVar(&bmcPass, "pass", viper.GetString("BMCPASS"),
		"BMC password")
	setCmd.Flags().StringVar(&bmcAddr, "addr", viper.GetString("BMCADDR"),
		"BMC IPv4 address")
}

// setCredentials updates a Credentials entity on Google Cloud Datastore.
func setCredentials() {
	bmcHost = makeBMCHostname(bmcHost, nameVersion)

	log.Infof("Updating credentials for host %v", bmcHost)
	provider, err := credsNewProvider(&creds.DatastoreConnector{}, projectID, namespace)
	rtx.Must(err, "Cannot connect to Datastore")
	defer provider.Close()

	creds, err := provider.FindCredentials(context.Background(), bmcHost)
	if err != nil {
		log.Errorf("Error while retrieving credentials for %s: %v", bmcHost, err)
		osExit(1)
	}

	// Update fields. When a new value is not specified, the corresponding
	// field is left as is.
	if bmcAddr != "" {
		creds.Address = bmcAddr
	}
	if bmcUser != "" {
		creds.Username = bmcUser
	}
	if bmcPass != "" {
		creds.Password = bmcPass
	}
	rtx.Must(provider.AddCredentials(context.Background(), bmcHost, creds),
		"Error while adding Credentials")
	fmt.Print(creds)
}
