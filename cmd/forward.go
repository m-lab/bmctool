package cmd

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	sshUser    string
	ports      []string
	tunnelHost string

	defaultPorts = []string{"4443:443", "5900"}
)

var forwardCmd = &cobra.Command{
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
func splitPorts(ports string) (int32, int32, error) {
	split := strings.Split(ports, ":")

	srcPort, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	if len(split) == 1 {
		return int32(srcPort), int32(srcPort), nil
	}

	dstPort, err := strconv.ParseInt(split[1], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	return int32(srcPort), int32(dstPort), nil
}

func forward(dstHost string) {
	if tunnelHost == "" || sshUser == "" {
		log.Error("BMCTUNNELHOST and BMCTUNNELUSER must not be empty.")
		osExit(1)
	}
	dstHost = makeBMCHostname(dstHost)

	errs, _ := errgroup.WithContext(context.Background())

	for _, port := range ports {
		srcPort, dstPort, err := splitPorts(port)
		if err != nil {
			log.Errorf("Cannot parse provided ports: %v", err)
			osExit(1)
		}

		log.Infof("Running ssh -N -q " + tunnelHost + " -L" +
			strconv.FormatInt(int64(srcPort), 10) + ":" + dstHost + ":" +
			strconv.FormatInt(int64(dstPort), 10))

		cmd := exec.Command("ssh", "-N", "-q", tunnelHost, "-L"+
			strconv.FormatInt(int64(srcPort), 10)+":"+dstHost+":"+
			strconv.FormatInt(int64(dstPort), 10))

		errs.Go(func() error {
			return cmd.Run()
		})
	}

	err := errs.Wait()
	if err != nil {
		log.Errorf("Error while running ssh command: %v", err)
	}
}
