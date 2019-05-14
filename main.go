package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

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
	// TODO: remove this once https://github.com/m-lab/reboot-service/issues/12
	// is closed.
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

func add() error {
	// Add a new BMC to Google Cloud Datastore
	var err error
	if *bmcUser == "" || *bmcPass == "" || *bmcAddr == "" {
		return errors.New("bmcuser, bmcpass and bmcaddr are required")
	}

	creds := &creds.Credentials{
		Address:  *bmcAddr,
		Hostname: *node,
		Model:    "DRAC",
		Username: *bmcUser,
		Password: *bmcPass,
	}

	log.Printf("Adding credentials for host %v\n", *node)
	provider := credsNewProvider(*projectID, namespace)

	// Provider.AddCredentials will create the entity regardless of whether it
	// exists already or not, so we need to explicitly check to prevent
	// overriding the existing entity by mistake.
	_, err = provider.FindCredentials(context.Background(), *node)
	if err == nil {
		return fmt.Errorf("Credentials for hostname %v already exist "+
			"- did you mean -update?", *node)
	}

	err = provider.AddCredentials(context.Background(), *node, creds)
	if err != nil {
		return err
	}

	return nil
}

func fetch() error {
	// Default behavior (if no other actions have been specified) is to fetch
	// credentials for the selected node.
	provider := credsNewProvider(*projectID, namespace)
	creds, err := provider.FindCredentials(context.Background(), *node)
	if err != nil {
		return err
	}

	jsonOutput, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonOutput))
	return nil
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
	//
	// -add: create a new entity
	// -update: updates an existing entity (TODO)
	// no flags: fetch credentials
	if *addAction {
		rtx.Must(add(), "Error while adding node")
	} else {
		rtx.Must(fetch(), "Error while fetching credentials")
	}
}
