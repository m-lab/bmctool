package cmd

import (
	"context"
	"testing"

	"github.com/m-lab/reboot-service/creds"
	"github.com/m-lab/reboot-service/creds/credstest"
	"github.com/stretchr/testify/assert"
)

func Test_setCredentials(t *testing.T) {
	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "127.0.0.1",
		Hostname: "mlab4d-lga0t.mlab-sandbox.measurement-lab.org",
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

	oldCredsNewProvider := credsNewProvider
	projectID = ""

	// Set up a FakeProvider with fake credentials.
	prov := credstest.NewProvider()
	prov.AddCredentials(context.Background(), "mlab4d-lga0t.mlab-sandbox.measurement-lab.org", fakeCreds)
	credsNewProvider = func(creds.Connector, string, string) (creds.Provider, error) {
		return prov, nil
	}

	// setCredentials should successfully change an existing entity
	bmcHost = "mlab4d-lga0t.mlab-sandbox.measurement-lab.org"
	bmcUser = "testuser"
	bmcPass = "testpass"
	bmcAddr = "127.0.0.2"
	setCredentials()

	// Check the node that's been just added.
	c, err := prov.FindCredentials(context.Background(),
		"mlab4d-lga0t.mlab-sandbox.measurement-lab.org")
	if err != nil {
		t.Errorf("FindCredentials() returned error: %v", err)
	}
	if c.Hostname != "mlab4d-lga0t.mlab-sandbox.measurement-lab.org" ||
		c.Username != "testuser" || c.Password != "testpass" ||
		c.Address != "127.0.0.2" || c.Model != "DRAC" {
		t.Errorf("setCredentials() didn't update the expected entity: %v", c)
	}

	// bmctool set should fail if called on a non-existing host.
	bmcHost = "mlab1-abc01.mlab-sandbox.measurement-lab.org"
	assert.PanicsWithValue(t, "os.Exit called", setCredentials,
		"os.Exit was not called")

	credsNewProvider = oldCredsNewProvider
}
