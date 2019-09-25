package forwarder

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/apex/log"
)

var osExit = os.Exit

type Forwarder interface {
	Start(context.Context) error
}

type Port struct {
	Src, Dst int
}

// sshForwarder allows to forward one or more ports via SSH tunneling.
// It relies on OpenSSH client, thus it needs to be installed on the system.
type sshForwarder struct {
	// Ports to forwards
	ports []Port

	// Tunnel host
	tHost string

	// Destination host
	dstHost string
}

func New(tHost string, dstHost string, ports []Port) Forwarder {
	return &sshForwarder{
		tHost:   tHost,
		dstHost: dstHost,
		ports:   ports,
	}
}

func (f *sshForwarder) getPortParams() []string {
	params := []string{}
	for _, p := range f.ports {
		params = append(params, fmt.Sprintf("-L%d:%s:%d", p.Src, f.dstHost, p.Dst))
	}
	return params
}

func (f *sshForwarder) Start(ctx context.Context) error {
	portParams := f.getPortParams()

	args := []string{"ssh", "-N", "-q", f.tHost}
	args = append(args, portParams...)
	log.Infof("Running %v", args)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	// monitor tunnel creation
	for end := time.Now().Add(5 * time.Second); time.Now().Before(end); {
		_, err = net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(f.ports[0].Src)), 1*time.Second)
		if err == nil {
			log.Info("SSH tunnel set up successfully! Press CTRL+C to exit.")
			return nil
		}
	}

	return err
}
