package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/creds"
	"github.com/spf13/cobra"
)

const (
	namespace        = "reboot-api"
	prodProjectID    = "mlab-oti"
	stagingProjectID = "mlab-staging"
	sandboxProjectID = "mlab-sandbox"
)

var (
	projectID string

	sandboxRegex = regexp.MustCompile("[a-zA-Z]{3}[0-9]t")
	stagingRegex = regexp.MustCompile("^mlab4")

	// These allow for testing.
	credsNewProvider = creds.NewProvider
	osExit           = os.Exit
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use: "bmctool",
		Long: `bmctool is a tool to manage BMCs on the M-Lab infrastructure.

It allows to read/create/update Credentials entities on Google Cloud Datastore,
set up an SSH tunnel (e.g. to connect through a trusted host), and reboot nodes.`,
	}
)

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
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "",
		"Project ID to use")
}

// parseNodeSite extracts node and site from a full hostname.
func parseNodeSite(hostname string) (string, string, error) {
	regex := regexp.MustCompile("(mlab[1-4]d?)\\.([a-zA-Z]{3}[0-9t]{2}).*")
	result := regex.FindStringSubmatch(hostname)
	if len(result) != 3 {
		return "", "",
			fmt.Errorf("The specified hostname is not a valid M-Lab node: %s", hostname)
	}

	return result[1], result[2], nil
}

// makeBMCHostname returns a full BMC hostname. There are different ways the
// hostname can be provided:
// - mlab1.lga0t
// - mlab1d.lga0t
// - mlab1.lga0t.measurement-lab.org
// - mlab1d.lga0t.measurement-lab.org
// This function returns the full hostname in any of these cases
func makeBMCHostname(name string) string {
	node, site, err := parseNodeSite(name)
	rtx.Must(err, "Cannot extract BMC hostname")

	// All the BMC hostnames must end with "d".
	if node[len(node)-1:] != "d" {
		node = node + "d"
	}
	return fmt.Sprintf("%s.%s.measurement-lab.org", node, site)
}

// getProjectID returns the correct GCP project to use based on the hostname.
func getProjectID(host string) string {
	if sandboxRegex.MatchString(host) {
		return sandboxProjectID
	}
	if stagingRegex.MatchString(host) {
		return stagingProjectID
	}

	return prodProjectID
}
