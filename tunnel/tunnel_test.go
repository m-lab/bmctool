package tunnel

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	sshserver "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
)

func TestSSHTunnel_Start(t *testing.T) {
	handlerFunc := func(s sshserver.Session) {
		io.WriteString(s, "test")
	}

	// Create intermediate SSH server.
	bounceSSHListener, err := net.Listen("tcp", ":3000")
	bounceSSHServer := &sshserver.Server{
		Handler: handlerFunc,
		LocalPortForwardingCallback: func(ctx sshserver.Context,
			destinationHost string, destinationPort uint32) bool {
			return true
		},
	}
	if err != nil {
		t.Errorf("Cannot create listener: %v", err)
	}
	go func() {
		log.Fatal(bounceSSHServer.Serve(bounceSSHListener))
	}()

	// Create destination SSH server.
	destSSHServer, err := net.Listen("tcp", "127.0.0.1:4000")
	if err != nil {
		t.Errorf("Cannot create listener: %v", err)
	}
	go func() {
		log.Fatal(sshserver.Serve(destSSHServer, handlerFunc))
	}()

	sshConfig := &ssh.ClientConfig{
		User: "test",
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	tun := &SSHTunnel{
		Local: &Endpoint{
			Host: "127.0.0.1",
			Port: 2000,
		},
		Server: &Endpoint{
			Host: "127.0.0.1",
			Port: int32(bounceSSHListener.Addr().(*net.TCPAddr).Port),
		},
		Remote: &Endpoint{
			Host: "127.0.0.1",
			Port: int32(destSSHServer.Addr().(*net.TCPAddr).Port),
		},
		Config: sshConfig,
	}

	tun.Start()

	<-context.Background().Done()

}
