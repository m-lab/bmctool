package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/apex/log"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/creds"
	"github.com/spf13/cobra"
)

var (
	bmcUser, bmcPass string
	bmcAddr, bmcHost string
	projectID        string

	// These allow for testing.
	credsNewProvider = creds.NewProvider
	osExit           = os.Exit
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <hostname> <address>",
	Short: "Add a new BMC",
	Long: `This command creates a new Credentials entity on GCD.

To use it, you also need to set the BMCUSER and BMCPASS environment variables
to an appropriate value.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		bmcHost = args[0]
		bmcAddr = args[1]

		rtx.Must(addCredentials(), "Error while adding credentials")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&bmcUser, "bmcuser", viper.GetString("BMCUSER"),
		"BMC username")
	addCmd.Flags().StringVar(&bmcPass, "bmcpass", viper.GetString("BMCPASS"),
		"BMC password")

	addCmd.MarkFlagRequired("bmcuser")
	addCmd.MarkFlagRequired("bmcpass")
}

// addCredentials adds a new BMC to Google Cloud Datastore.
func addCredentials() error {
	if bmcUser == "" || bmcPass == "" || bmcAddr == "" {
		log.Error("bmcuser, bmcpassword and addr are required")
		osExit(1)
	}

	creds := &creds.Credentials{
		Address:  bmcAddr,
		Hostname: bmcHost,
		Model:    "DRAC",
		Username: bmcUser,
		Password: bmcPass,
	}

	log.Infof("Adding credentials for host %v", bmcHost)
	provider := credsNewProvider(projectID, namespace)

	// Provider.AddCredentials will create the entity regardless of whether it
	// exists already or not, so we need to explicitly check to prevent
	// overriding the existing entity by mistake.
	_, err := provider.FindCredentials(context.Background(), bmcHost)
	if err == nil {
		log.Errorf("Credentials for hostname %v already exist", bmcHost)
		osExit(1)
	}

	rtx.Must(provider.AddCredentials(context.Background(), bmcHost, creds),
		"Error while adding Credentials")

	fmt.Print(creds)
	return nil
}
