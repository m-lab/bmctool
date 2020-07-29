package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/m-lab/bmctool/forwarder"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/reboot-service/connector"
	"github.com/m-lab/reboot-service/creds"
	"github.com/spf13/cobra"
)

const (
	adminIdx = 2 // Index for the 'admin' user on DRACs
)

var (
	bmcPort   int32
	localPort int32
	useTunnel bool
	// keysSetCmd represents the keys set command
	keysSetCmd = &cobra.Command{
		Use:   "set <host> <index> <key>",
		Short: "Replaces the SSH key at <index> with <key>",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			setKey(args[0], args[1], args[2])
		},
	}
)

func init() {
	keysCmd.AddCommand(keysSetCmd)

	keysSetCmd.Flags().Int32Var(&bmcPort, "bmcport", defaultBMCPort,
		"BMC port to use")
	keysSetCmd.Flags().BoolVar(&useTunnel, "tunnel", false,
		"Tunnel through an intermediate host")
	keysSetCmd.Flags().Int32Var(&localPort, "localport", defaultLocalPort,
		"Local port to use when tunneling with -tunnel")
}

func setKey(host, idx, key string) {
	bmcNode := makeBMCHostname(host)

	log.Infof("Project: %s", bmcNode.Project)
	log.Infof("Fetching credentials for %s", bmcNode.String())

	provider, err := credsNewProvider(&creds.DatastoreConnector{}, bmcNode.Project, namespace)
	rtx.Must(err, "Cannot connect to Datastore")
	defer provider.Close()

	creds, err := provider.FindCredentials(context.Background(), bmcNode.String())
	rtx.Must(err, "Cannot fetch credentials")

	// Make a connection to the host
	connectionConfig := &connector.ConnectionConfig{
		Hostname: bmcNode.String(),
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
		sshForwarder := forwarder.New(tunnelHost, sshUser, bmcNode.String(), ports)
		sshForwarder.Start(context.Background())
		connectionConfig.Hostname = "127.0.0.1"
		connectionConfig.Port = localPort
	}

	conn, err := connector.NewConnector().NewConnection(connectionConfig)
	rtx.Must(err, "Cannot connect to BMC: %s", bmcNode.String())
	defer conn.Close()

	// Sending the racadm command via SSH in single-command mode means the
	// SSH key will be truncated. Apparently, the only way to make this work
	// is to request a shell and run the command interactively, then check
	// stdout/stderr for signs that the command execution succeeded.
	cmd := fmt.Sprintf("racadm sshpkauth -i %d -k %s -t \"%s\"", adminIdx, idx, key)
	log.Infof("Running command: %s", cmd)
	out, err := conn.ExecDRACShell(cmd)
	rtx.Must(err, "Cannot set SSH key on %s (index: %s): %s", bmcNode.String(), idx, out)

	if !strings.Contains(out, "PK SSH Authentication operation completed successfully.") {
		log.Errorf("Operation failed: %s", out)
		osExit(1)
	}

	log.Info("The SSH key has been added successfully")
}
