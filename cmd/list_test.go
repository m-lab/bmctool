package cmd

import (
	"context"
	"testing"

	"github.com/m-lab/reboot-service/creds"
	"github.com/m-lab/reboot-service/creds/credstest"
	"github.com/stretchr/testify/assert"
)

func Test_listBMCs(t *testing.T) {
	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "127.0.0.1",
		Hostname: "mlab4d.lga0t.measurement-lab.org",
		Username: "username",
		Password: "password",
		Model:    "DRAC",
	}

	oldCredsNewProvider := credsNewProvider

	// Set up a FakeProvider with fake credentials.
	prov := credstest.NewProvider()
	prov.AddCredentials(context.Background(), "mlab4d.lga0t.measurement-lab.org", fakeCreds)
	credsNewProvider = func(creds.Connector, string, string) (creds.Provider, error) {
		return prov, nil
	}

	// Replace osExit so that tests don't stop running.
	osExit = func(code int) {
		if code != 1 {
			t.Fatalf("Expected a 1 exit code, got %d.", code)
		}

		panic("os.Exit called")
	}

	// listBMCs() should fails when called without --project.
	projectID = ""
	assert.PanicsWithValue(t, "os.Exit called", listBMCs,
		"os.Exit was not called")

	projectID = "test-project"
	listBMCs()
	credsNewProvider = oldCredsNewProvider

}
