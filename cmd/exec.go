package cmd

import (
	"context"

	"github.com/apex/log"
	"github.com/m-lab/bmctool/forwarder"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/connector"
	"github.com/m-lab/reboot-service/creds"
	"github.com/spf13/cobra"
)

var newConnector = connector.NewConnector

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec <host> <command>",
	Short: "Execute a single command on a BMC",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		toExec := args[1]
		exec(host, toExec)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().Int32Var(&bmcPort, "bmcport", defaultBMCPort,
		"BMC port to use")
	execCmd.Flags().Int32Var(&localPort, "localport", defaultLocalPort,
		"Local port to use when tunneling with -tunnel")
	execCmd.Flags().BoolVar(&useTunnel, "tunnel", false,
		"Tunnel through an intermediate host")
}

func exec(host, cmd string) {
	bmcHost := makeBMCHostname(host)

	log.Infof("Project: %s", projectID)
	log.Infof("Fetching credentials for %s", bmcHost)
	provider, err := credsNewProvider(&creds.DatastoreConnector{}, projectID, namespace)
	rtx.Must(err, "Cannot connect to Datastore")
	defer provider.Close()

	creds, err := provider.FindCredentials(context.Background(), bmcHost)
	rtx.Must(err, "Cannot fetch credentials")

	// Make a connection to the BMC.
	connectionConfig := &connector.ConnectionConfig{
		Hostname: bmcHost,
		Username: creds.Username,
		Password: creds.Password,
		Port:     bmcPort,
		ConnType: connector.BMCConnection,
		Timeout:  bmcTimeout,
	}

	if useTunnel {
		if tunnelHost == "" || sshUser == "" {
			log.Error("BMCTUNNELHOST and BMCTUNNELUSER must not be empty.")
			osExit(1)
		}
		ports := []forwarder.Port{
			{
				Src: int(localPort),
				Dst: int(bmcPort),
			},
		}
		sshForwarder := newForwarder(tunnelHost, sshUser, bmcHost, ports)
		sshForwarder.Start(context.Background())
		connectionConfig.Hostname = "127.0.0.1"
		connectionConfig.Port = localPort
	}

	// Establish connection to the BMC.
	c := newConnector()
	conn, err := c.NewConnection(connectionConfig)
	rtx.Must(err, "Cannot connect to BMC: %s", bmcHost)
	defer conn.Close()

	// Execute the command.
	log.Infof("Running command: %s", cmd)
	out, err := conn.ExecDRACShell(cmd)
	rtx.Must(err, "Cannot execute command \"%s\"", cmd)

	log.Info(out)
}
