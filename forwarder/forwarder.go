package forwarder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/apex/log"
)

var osExit = os.Exit

// SSHForwarder allows to forward one or more ports via SSH tunneling.
// It relies on OpenSSH client, thus it needs to be installed on the system.
type SSHForwarder struct {
	// Ports to forwards ("src:dst" or "dst")
	ports []string

	// Tunnel host/port
	tHost string

	// Destination host
	dstHost string
}

func NewSSHForwarder(tHost string, dstHost string, ports []string) *SSHForwarder {
	return &SSHForwarder{
		tHost:   tHost,
		dstHost: dstHost,
		ports:   ports,
	}
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

func (f *SSHForwarder) getPortParams() ([]string, error) {
	params := []string{}
	for _, port := range f.ports {
		srcPort, dstPort, err := splitPorts(port)
		if err != nil {
			return nil, err
		}
		params = append(params, fmt.Sprintf("-L%d:%s:%d", srcPort, f.dstHost, dstPort))
	}
	return params, nil
}

func (f *SSHForwarder) Start(ctx context.Context) error {
	portParams, err := f.getPortParams()
	if err != nil {
		return err
	}

	args := []string{"ssh", "-N", "-q", f.tHost}
	args = append(args, portParams...)
	log.Infof("Running %v", args)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	return cmd.Start()
}
