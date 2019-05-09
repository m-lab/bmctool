package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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

func main() {
	flag.Usage = usage
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Could not parse env vars")

	// No node specified, nothing to do.
	if *node == "" {
		fmt.Fprintln(os.Stderr, "Error: node not specified.")
		flag.Usage()
	}

	provider := credsNewProvider(*projectID, namespace)
	creds, err := provider.FindCredentials(context.Background(), *node)
	rtx.Must(err, "Error while fetching credentials: %v\n", err)

	jsonOutput, err := json.MarshalIndent(creds, "", "  ")
	rtx.Must(err, "Cannot marshal JSON: %v\n")

	fmt.Println(string(jsonOutput))
}
