package cmd

import (
	"context"
	"testing"

	"github.com/m-lab/reboot-service/creds"
	"github.com/m-lab/reboot-service/creds/credstest"
)

func Test_printCredentials(t *testing.T) {
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
	credsNewProvider = func(string, string) creds.Provider {
		return prov
	}

	// printCredentials is intentionally called with a short name here.
	printCredentials("mlab4.lga0t")

	credsNewProvider = oldCredsNewProvider
}
