package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"

	"github.com/apex/log"
	"github.com/m-lab/reboot-service/creds"
)

const (
	defaultProjectID = "mlab-sandbox"
	namespace        = "reboot-api"
)

var (
	node      = flag.String("node", "", "The node's name.")
	projectID = flag.String("project", defaultProjectID, "Project ID to use.")

	// These allows for testing.
	createProvider = creds.NewProvider
	marshalJSON    = json.MarshalIndent
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "%s -node <node>: fetch credentials for <node>.\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Could not parse env vars")

	// No node specified, nothing to do.
	if *node == "" {
		fmt.Fprintln(os.Stderr, "Error: node not specified.")
		flag.Usage()
		return
	}

	provider := createProvider(*projectID, namespace)
	creds, err := provider.FindCredentials(context.Background(), *node)
	if err != nil {
		log.Errorf("Error while fetching credentials: %v\n", err)
		return
	}

	jsonOutput, err := marshalJSON(creds, "", "  ")
	if err != nil {
		log.Errorf("Cannot marshall JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonOutput))
}
