package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"

	"github.com/m-lab/reboot-service/creds"
)

const (
	defaultProjectID = "mlab-sandbox"
	namespace        = "reboot-api"
)

var (
	node      = flag.String("node", "", "The node's name.")
	projectID = flag.String("project", defaultProjectID, "Project ID to use.")
	bmcUser   = flag.String("bmcuser", "",
		"BMC username (for adding or updating a BMC.)")
	bmcPass = flag.String("bmcpassword", "",
		"BMC password (for adding or updating a BMC.)")
	bmcAddr = flag.String("addr", "",
		"BMC IP address (for adding or updating a BMC.)")

	// Actions
	addAction = flag.Bool("add", false, "Add a new node to GCS.")

	// These allow for testing.
	credsNewProvider = creds.NewProvider
	osExit           = os.Exit
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "%s -node <node>: fetch credentials for <node>.\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	osExit(1)
}

// addCredentials adds a new BMC to Google Cloud Datastore.
func addCredentials() error {
	if *bmcUser == "" || *bmcPass == "" || *bmcAddr == "" {
		log.Error("bmcuser, bmcpassword and addr are required")
		osExit(1)
	}

	creds := &creds.Credentials{
		Address:  *bmcAddr,
		Hostname: *node,
		Model:    "DRAC",
		Username: *bmcUser,
		Password: *bmcPass,
	}

	log.Infof("Adding credentials for host %v", *node)
	provider := credsNewProvider(*projectID, namespace)

	// Provider.AddCredentials will create the entity regardless of whether it
	// exists already or not, so we need to explicitly check to prevent
	// overriding the existing entity by mistake.
	_, err := provider.FindCredentials(context.Background(), *node)
	if err == nil {
		log.Errorf("Credentials for hostname %v already exist", *node)
		osExit(1)
	}

	rtx.Must(provider.AddCredentials(context.Background(), *node, creds),
		"Error while adding Credentials")

	fmt.Print(creds)
	return nil
}

// printCredentials retrieves credentials for a given hostname and prints them
// in JSON format.
func printCredentials(host string) {
	provider := credsNewProvider(*projectID, namespace)
	creds, err := provider.FindCredentials(context.Background(), *node)
	rtx.Must(err, "Cannot fetch credentials")

	fmt.Print(creds)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Could not parse env vars")

	// No node specified, nothing to do.
	if *node == "" {
		fmt.Fprintln(os.Stderr, "Error: node not specified.")
		flag.Usage()
	}

	// Handle action flags
	// TODO(roberto): use a package that handles subcommands (such as kingpin)
	// instead of "flag".
	//
	// -add: create a new entity
	// -update: updates an existing entity (TODO)
	// no flags: fetch credentials
	if *addAction {
		rtx.Must(addCredentials(), "Error while adding credentials")
	} else {
		// Default behavior (if no other actions have been specified) is to fetch
		// credentials for the requested node.
		printCredentials(*node)
	}
}
