package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	namespace        = "reboot-api"
	defaultProjectID = "mlab-sandbox"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "bmctool",
	Long: `bmctool is a tool to manage BMCs on the M-Lab infrastructure.

It allows to read/create/update Credentials entities on Google Cloud Datastore,
set up an SSH tunnel (e.g. to connect through a trusted host), and reboot nodes.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// The --project flag is used by several commands, thus it's defined
	// as global ("Persistent") flag here.
	addCmd.PersistentFlags().StringVar(&projectID, "project", defaultProjectID,
		"Project ID to use")
}
