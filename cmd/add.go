package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/apex/log"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/creds"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var (
	bmcUser, bmcPass string
	bmcHost, bmcAddr string
	addCmd           = &cobra.Command{
		Use:   "add <hostname> <address>",
		Short: "Add a new BMC",
		Long: `This command creates a new Credentials entity on GCD.

To use it, you also need to set the BMCUSER and BMCPASS environment variables
to an appropriate value.`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			bmcHost = args[0]
			bmcAddr = args[1]

			addCredentials()
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)

	viper.AutomaticEnv()

	addCmd.Flags().StringVar(&bmcUser, "bmcuser", viper.GetString("BMCUSER"),
		"BMC username")
	addCmd.Flags().StringVar(&bmcPass, "bmcpass", viper.GetString("BMCPASS"),
		"BMC password")

}

// addCredentials adds a new BMC to Google Cloud Datastore.
func addCredentials() {
	if bmcUser == "" || bmcPass == "" {
		log.Error("BMCUSER and BMCPASS must not be empty.")
		osExit(1)
	}

	bmcHost = makeBMCHostname(bmcHost)

	c := &creds.Credentials{
		Address:  bmcAddr,
		Hostname: bmcHost,
		Model:    "DRAC",
		Username: bmcUser,
		Password: bmcPass,
	}

	log.Infof("Adding credentials for host %v", bmcHost)

	provider, err := credsNewProvider(&creds.DatastoreConnector{}, projectID, namespace)
	rtx.Must(err, "Cannot connect to Datastore")
	defer provider.Close()

	// Provider.AddCredentials will create the entity regardless of whether it
	// exists already or not, so we need to explicitly check to prevent
	// overriding the existing entity by mistake.
	_, err = provider.FindCredentials(context.Background(), bmcHost)
	if err == nil {
		log.Errorf("Credentials for hostname %v already exist", bmcHost)
		osExit(1)
	}

	rtx.Must(provider.AddCredentials(context.Background(), bmcHost, c),
		"Error while adding Credentials")

	fmt.Print(c)
}
