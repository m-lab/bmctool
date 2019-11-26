package cmd

import (
	"context"
	"testing"

	"github.com/m-lab/reboot-service/connector"

	"github.com/m-lab/reboot-service/creds"
	"github.com/m-lab/reboot-service/creds/credstest"
)

// Mock objects for Connector/Connection.
type mockConnector struct{}

type mockConnection struct{}

func (connection *mockConnection) ExecDRACShell(string) (string, error) {
	return "Not implemented.", nil
}

func (connection *mockConnection) Reboot() (string, error) {
	return "Not implemented.", nil
}
func (connection *mockConnection) Close() error {
	return nil
}

func (connector *mockConnector) NewConnection(config *connector.ConnectionConfig) (connector.Connection, error) {
	return &mockConnection{}, nil
}

func Test_exec(t *testing.T) {
	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "127.0.0.1",
		Hostname: "mlab1d.tst01.lga0t.measurement-lab.org",
		Username: "username",
		Password: "password",
		Model:    "DRAC",
	}

	oldCredsNewProvider := credsNewProvider
	oldNewConnector := newConnector
	oldNewForwarder := newForwarder

	// Set up a FakeProvider with fake credentials.
	prov := credstest.NewProvider()
	prov.AddCredentials(context.Background(), "mlab1d.tst01.measurement-lab.org", fakeCreds)
	credsNewProvider = func(creds.Connector, string, string) (creds.Provider, error) {
		return prov, nil
	}

	newConnector = func() connector.Connector {
		return &mockConnector{}
	}

	newForwarder = newForwarderMock

	useTunnel = true
	tunnelHost = "test"
	bmcUser = "test"
	exec("mlab1d.tst01", "help")

	newForwarder = oldNewForwarder
	credsNewProvider = oldCredsNewProvider
	newConnector = oldNewConnector
}
