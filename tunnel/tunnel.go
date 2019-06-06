package tunnel

import (
	"io"
	"net"
	"os"

	"github.com/apex/log"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHTunnel represents an SSH tunnel.
type SSHTunnel struct {
	// Local server endpoint
	Local *Endpoint

	// Intermediate server endpoint
	Server *Endpoint

	// Remote server endpoint
	Remote *Endpoint

	// Client configuration
	Config *ssh.ClientConfig
}

// Start initializes the SSH tunnel.
func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		log.Errorf("Cannot listen on %s: %v", tunnel.Local, err)
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Cannot accept connection: %v", err)
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		log.Errorf("Server dial error: %s", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		log.Errorf("Remote dial error: %s", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Debugf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

// SSHAgent gets a ssh.AuthMethod from the local ssh-agent instance (if any).
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
