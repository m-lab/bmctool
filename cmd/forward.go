package cmd

import (
	"context"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/m-lab/bmctool/forwarder"
	"github.com/m-lab/go/rtx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sshUser    string
	ports      []string
	tunnelHost string

	defaultPorts = []string{"4443:443", "5900"}
)

var (
	forwardCmd = &cobra.Command{
		Use:   "forward <host>",
		Short: "Forward ports via an SSH tunnel",
		Long: `This command creates an SSH tunnel to a given <host>.

Ports to be forwarded can be specified with the (repeatable) --port flag.
Local and remote ports can be specified with the following syntax:

--port src[:dest]

e.g.:

bmctool forward <host> --port 4443:443 --port 5900

If dest is unspecified, it'll be the same as src.

The host to use for tunneling can be specified via the --tunnel-host flag,
or the BMCTUNNELHOST environment variable.

The username to use to connect to the intermediate host can be specified via
the --username flag, or the BMCTUNNELUSER environment variable.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dstHost := args[0]
			forward(dstHost)
		},
	}
)

func init() {
	rootCmd.AddCommand(forwardCmd)

	viper.AutomaticEnv()

	forwardCmd.Flags().StringArrayVar(&ports, "port", defaultPorts, "source:destination")
	forwardCmd.Flags().StringVar(&tunnelHost, "tunnel-host",
		viper.GetString("BMCTUNNELHOST"), "intermediate host")
	forwardCmd.Flags().StringVar(&sshUser, "username",
		viper.GetString("BMCTUNNELUSER"), "username for intermediate host")

}

// splitPorts takes a string containing either a "local:remote" ports pair
// or just "port" and returns local/remote as separate variables. If the string
// contains a single port, it returns the same port for local and remote.
func splitPorts(ports string) (forwarder.Port, error) {
	split := strings.Split(ports, ":")

	srcPort, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return forwarder.Port{}, err
	}

	if len(split) == 1 {
		return forwarder.Port{Src: int(srcPort), Dst: int(srcPort)}, nil
	}

	dstPort, err := strconv.ParseInt(split[1], 10, 32)
	if err != nil {
		return forwarder.Port{}, err
	}

	return forwarder.Port{Src: int(srcPort), Dst: int(dstPort)}, nil
}

func forward(dstHost string) {
	if tunnelHost == "" || sshUser == "" {
		log.Error("BMCTUNNELHOST and BMCTUNNELUSER must not be empty.")
		osExit(1)
	}
	dstHost = makeBMCHostname(dstHost, nameVersion)

	portFwd := []forwarder.Port{}
	for _, port := range ports {
		p, err := splitPorts(port)
		rtx.Must(err, "Cannot parse provided port")
		portFwd = append(portFwd, p)
	}
	forwarder := newForwarder(tunnelHost, sshUser, dstHost, portFwd)

	ctx := context.Background()
	rtx.Must(forwarder.Start(context.Background()), "Cannot start SSH tunnel")

	<-ctx.Done()
}
