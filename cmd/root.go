package cmd

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/m-lab/bmctool/forwarder"
	"github.com/m-lab/go/host"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/creds"
	"github.com/spf13/cobra"
)

const (
	namespace        = "reboot-api"
	prodProjectID    = "mlab-oti"
	stagingProjectID = "mlab-staging"
	sandboxProjectID = "mlab-sandbox"
	defaultBMCPort   = 806
	defaultLocalPort = 8060
	bmcTimeout       = 30 * time.Second
)

var (
	projectID   string
	nameVersion string

	// TODO(kinkade): these patterns should go away in favor of determining the
	// project from siteinfo.
	sandboxRegex = regexp.MustCompile("[a-zA-Z]{3}[0-9]t")
	stagingRegex = regexp.MustCompile("^mlab4")

	// These allow for testing.
	credsNewProvider = creds.NewProvider
	osExit           = os.Exit
	newForwarder     = forwarder.New
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
	// The --project flags is used by several commands, thus it's defined
	// as global ("Persistent") flags here.
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "",
		"Project ID to use")
}

// parseNodeSite extracts node and site from a full hostname.
func parseNodeSite(hostname string) (string, string, error) {
	regex := regexp.MustCompile(`(mlab[1-4]d?)[.-]([a-zA-Z]{3}[0-9ct]{2}).*`)
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
// - mlab1-lga0t
// - mlab1d.lga0t
// - mlab1d-lga0t
// - mlab1.lga0t.measurement-lab.org
// - mlab1-lga0t.measurement-lab.org
// - mlab1d.lga0t.measurement-lab.org
// - mlab1d-lga0t.measurement-lab.org
// This function returns the full hostname in any of these cases
func makeBMCHostname(name string) string {
	node, site, err := parseNodeSite(name)
	rtx.Must(err, "Cannot extract BMC hostname")

	if projectID == "" {
		projectID = getProjectID(name)
	}

	// All the BMC hostnames must end with "d".
	if node[len(node)-1:] != "d" {
		node = node + "d"
	}

	return fmt.Sprintf("%s-%s.%s.measurement-lab.org", node, site, projectID)
}

// getProjectID returns the correct GCP project to use based on the hostname.
func getProjectID(hostname string) string {
	// First, try parsing the hostname with host.Parse().
	parsed, err := host.Parse(hostname)
	if err == nil && parsed.Project != "" {
		return parsed.Project
	}
	// If host.Parse() fails, try with regular expressions.
	// TODO: replace this with siteinfo's projects.json.
	if sandboxRegex.MatchString(hostname) {
		return sandboxProjectID
	}
	if stagingRegex.MatchString(hostname) {
		return stagingProjectID
	}

	return prodProjectID
}
