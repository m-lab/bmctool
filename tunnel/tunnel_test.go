package tunnel

import (
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	sshserver "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
)

func TestSSHTunnel_Start(t *testing.T) {
	handlerFunc := func(s sshserver.Session) {
		io.WriteString(s, "test")
	}

	// Create intermediate SSH server.
	bounceSSHListener, err := net.Listen("tcp", ":0")
	bounceSSHServer := &sshserver.Server{
		Handler: handlerFunc,
		LocalPortForwardingCallback: func(ctx sshserver.Context,
			destinationHost string, destinationPort uint32) bool {
			return true
		},
	}
	if err != nil {
		t.Fatalf("Cannot create listener: %v", err)
	}
	go func() {
		log.Fatal(bounceSSHServer.Serve(bounceSSHListener))
	}()

	// Create destination SSH server.
	destSSHServer, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Cannot create listener: %v", err)
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
			Port: int32(destSSHServer.Addr().(*net.TCPAddr).Port) + 1,
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

	go func() { tun.Start() }()

	time.Sleep(2 * time.Second)

	// Connect to the tunnel and verify that the received message is the
	// expected one from the remote server.
	cl, err := ssh.Dial("tcp", tun.Local.String(), sshConfig)
	if err != nil {
		t.Fatalf("Cannot connect to the local endpoint: %v", err)
	}

	sess, err := cl.NewSession()
	if err != nil {
		t.Fatalf("Cannot create SSH session: %v", err)
	}

	sshout, err := sess.StdoutPipe()
	if err != nil {
		t.Fatalf("Cannot pipe stdout: %v", err)
	}

	err = sess.Shell()
	if err != nil {
		t.Fatalf("Cannot start shell: %v", err)
	}

	output := readBuffForString(sshout)
	if output != "test" {
		t.Fatalf("Unexpected output: %s", output)
	}
	fmt.Println("Done.")

}

func readBuffForString(sshOut io.Reader) string {
	buf := make([]byte, 1000)
	n, err := sshOut.Read(buf) //this reads the ssh terminal
	waitingString := ""
	if err == nil {
		waitingString = string(buf[:n])
	}
	return waitingString
}
