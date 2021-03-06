package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/m-lab/bmctool/forwarder"
	"github.com/m-lab/go/host"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/creds"
	"github.com/segmentio/encoding/json"
	"github.com/spf13/cobra"
)

const (
	namespace        = "reboot-api"
	defaultBMCPort   = 22
	defaultLocalPort = 8060
	bmcTimeout       = 30 * time.Second
	siteinfoBaseURL  = "https://siteinfo.mlab-oti.measurementlab.net/v2/"
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

	httpGet = http.Get
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
	// The --project flag is used by several commands, thus it is defined
	// as a global ("Persistent") flag here.
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "",
		"Project ID to use")
}

// parseNodeSite extracts node and site from a hostname.
func parseNodeSite(hostname string) (string, string, error) {
	regex := regexp.MustCompile(`(mlab[1-4]d?)[.-]([a-zA-Z]{3}[0-9ct]{2}).*`)
	result := regex.FindStringSubmatch(hostname)
	if len(result) != 3 {
		return "", "",
			fmt.Errorf("The specified hostname is not a valid M-Lab node: %s",
				hostname)
	}

	return result[1], result[2], nil
}

// makeBMCHostname returns a host.Name containing the components of the passed
// hostname. The accepted formats are:
// - mlab1.lga0t
// - mlab1-lga0t
// - mlab1d.lga0t
// - mlab1d-lga0t
// - mlab1.lga0t.measurement-lab.org
// - mlab1-lga0t.measurement-lab.org
// - mlab1d.lga0t.measurement-lab.org
// - mlab1d-lga0t.measurement-lab.org
// - mlab1-lga0t.mlab-sandbox.measurement-lab.org
// - mlab1d-lga0t.mlab-sandbox.measurement-lab.org
//
// If the hostname is not complete, this function retrieves the project ID
// from siteinfo. If the machine/site combination cannot be found in siteinfo,
// it returns an error.
func makeBMCHostname(name string) host.Name {
	// Is it a full M-Lab hostname?
	node, err := host.Parse(name)
	if err != nil {
		machine, site, err := parseNodeSite(name)
		rtx.Must(err, "Cannot extract BMC hostname")

		node = host.Name{
			Machine: machine,
			Site:    site,
			Domain:  "measurement-lab.org",
			Version: "v2",
		}
	}

	// Allow for manual overriding of the project ID.
	if projectID != "" {
		node.Project = projectID
	}

	// If the projectID was not specified and can't be inferred from the
	// hostname, get it from siteinfo.
	if node.Project == "" {
		// Make sure machine name does not end in 'd' already
		if strings.HasSuffix(node.Machine, "d") {
			node.Machine = node.Machine[:len(node.Machine)-1]
		}
		project, err := getProjectID(node)
		rtx.Must(err, "cannot get project ID from siteinfo for %s-%s",
			node.Machine, node.Site)
		node.Project = project
	}

	// All the BMC hostnames must end with "d".
	if !strings.HasSuffix(node.Machine, "d") {
		node.Machine = node.Machine + "d"
	}

	return node
}

// getProjectID returns the correct GCP project for a given node by looking up
// the sites/projects.json file from Siteinfo.
func getProjectID(node host.Name) (string, error) {
	// TODO(roberto) use m-lab/go/siteinfo once it supports the projects.json
	// output format.
	resp, err := httpGet(siteinfoBaseURL + "sites/projects.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var projects map[string]string
	rtx.Must(json.Unmarshal(body, &projects), "cannot parse projects.json")

	if project, ok := projects[fmt.Sprintf("%s-%s", node.Machine,
		node.Site)]; ok {
		return project, nil
	}

	return "", errors.New("hostname not found in projects.json")
}
