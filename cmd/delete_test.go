package cmd

import (
	"context"
	"testing"

	"github.com/m-lab/reboot-service/creds"
	"github.com/m-lab/reboot-service/creds/credstest"
	"github.com/stretchr/testify/assert"
)

func Test_deleteCredentials(t *testing.T) {
	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "127.0.0.1",
		Hostname: "mlab4d.lga0t.measurement-lab.org",
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

	// deleteCredentials should successfully remove an existing entity.
	bmcHost = "mlab4d-lga0t"
	deleteCredentials()

	_, err := prov.FindCredentials(context.Background(), "mlab4.lga0t")
	if err == nil {
		t.Errorf("deleteCredentials() did not delete Credentials.")
	}

	// deleteCredentials() should fail if provider.DeleteCredentials() fails.
	// Since our FakeProvider fails if trying to delete a non-existing entity,
	// calling deleteCredentials() again on the same key is enough to test it.
	assert.PanicsWithValue(t, "os.Exit called", deleteCredentials,
		"os.Exit was not called")

	credsNewProvider = oldCredsNewProvider
}
