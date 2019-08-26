package cmd

import (
	"context"
	"testing"

	"github.com/m-lab/reboot-service/creds"
	"github.com/m-lab/reboot-service/creds/credstest"
	"github.com/stretchr/testify/assert"
)

func Test_addCredentials(t *testing.T) {
	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "127.0.0.1",
		Hostname: "mlab4.lga0t.measurement-lab.org",
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

	// Set up a FakeProvider with fake credentials.
	prov := credstest.NewProvider()
	prov.AddCredentials(context.Background(), "mlab4d.lga0t.measurement-lab.org", fakeCreds)
	credsNewProvider = func(creds.Connector, string, string) (creds.Provider, error) {
		return prov, nil
	}

	// addCredentials should successfully add a new node.
	// The hostname is intentionally provided with the short name here.
	bmcHost = "mlab1.lga0t"
	bmcAddr = "127.0.0.1"
	bmcUser = "user"
	bmcPass = "pass"
	addCredentials()

	// Check the node that's been just added.
	c, err := prov.FindCredentials(context.Background(), "mlab1d.lga0t.measurement-lab.org")
	if err != nil {
		t.Errorf("FindCredentials() returned error: %v", err)
	}
	if c.Hostname != "mlab1d.lga0t.measurement-lab.org" ||
		c.Username != "user" || c.Password != "pass" ||
		c.Address != "127.0.0.1" || c.Model != "DRAC" {
		t.Errorf("AddCredentials() didn't add the expected entity: %v", c)
	}

	// bmctool add should fail if trying to add the same node again.
	assert.PanicsWithValue(t, "os.Exit called", addCredentials,
		"os.Exit was not called")

	// bmctool add should fail if bmcUser or bmcPass are empty.
	bmcUser = ""
	bmcPass = ""
	assert.PanicsWithValue(t, "os.Exit called", addCredentials,
		"os.Exit was not called")

	credsNewProvider = oldCredsNewProvider
}
