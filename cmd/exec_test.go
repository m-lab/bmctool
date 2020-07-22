package cmd

import (
	"context"
	"os"
	"testing"

	"github.com/m-lab/reboot-service/connector"

	"github.com/m-lab/reboot-service/creds"
	"github.com/m-lab/reboot-service/creds/credstest"
)

// Mock objects for Connector/Connection.
type mockConnector struct {
	conn *mockConnection
}

type mockConnection struct {
	execCalls int
}

func (connection *mockConnection) ExecDRACShell(string) (string, error) {
	connection.execCalls++
	return "ExecDRACShell called.", nil
}

func (connection *mockConnection) Reboot() (string, error) {
	return "Not implemented.", nil
}
func (connection *mockConnection) Close() error {
	return nil
}

func (connector *mockConnector) NewConnection(config *connector.ConnectionConfig) (connector.Connection, error) {
	connector.conn = &mockConnection{}
	return connector.conn, nil
}

func Test_exec(t *testing.T) {
	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "127.0.0.1",
		Hostname: "mlab1d-lga0t.mlab-sandbox.measurement-lab.org",
		Username: "username",
		Password: "password",
		Model:    "DRAC",
	}

	// Replace osExit so that tests don't stop running.
	osExit = func(code int) {
		if code != 1 {
			t.Fatalf("Expected a 1 exit code, got %d.", code)
		}

		panic("os.Exit called")
	}

	defer func() {
		osExit = os.Exit
	}()

	// Set up a FakeProvider with fake credentials.
	prov := credstest.NewProvider()
	prov.AddCredentials(context.Background(),
		"mlab1d-lga0t.mlab-sandbox.measurement-lab.org", fakeCreds)

	oldCredsNewProvider := credsNewProvider
	oldNewConnector := newConnector
	oldNewForwarder := newForwarder

	credsNewProvider = func(creds.Connector, string, string) (creds.Provider, error) {
		return prov, nil
	}

	c := &mockConnector{}
	newConnector = func() connector.Connector {
		return c
	}

	newForwarder = newForwarderMock

	useTunnel = true
	tunnelHost = "test"
	sshUser = "test"
	exec("mlab1d-lga0t", "help")

	if c.conn.execCalls != 1 {
		t.Errorf("exec called but execCalls != 1")
	}

	newForwarder = oldNewForwarder
	credsNewProvider = oldCredsNewProvider
	newConnector = oldNewConnector
}
