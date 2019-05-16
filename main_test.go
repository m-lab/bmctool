package main

import (
	"context"
	"testing"

	"github.com/m-lab/reboot-service/creds/credstest"

	"github.com/m-lab/go/osx"
	"github.com/m-lab/reboot-service/creds"
	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "0.0.0.0",
		Hostname: "mlab4.lga0t.measurement-lab.org",
		Username: "username",
		Password: "password",
		Model:    "drac",
	}

	// Replace osExit so that tests don't stop running.
	osExit = func(code int) {
		if code != 1 {
			t.Fatalf("Expected a 1 exit code, got %d.", code)
		}

		panic("os.Exit called")
	}

	oldCredsNewProvider := credsNewProvider

	// No node specified, just print usage and return.
	t.Run("failure-no-node-provided", func(t *testing.T) {
		assert.PanicsWithValue(t, "os.Exit called", main, "os.Exit was not called")
	})

	// Set up a FakeProvider with fake credentials.
	credsNewProvider = func(string, string) creds.Provider {
		prov := credstest.NewProvider()
		prov.AddCredentials(context.Background(), "mlab4d.lga0t.measurement-lab.org", fakeCreds)
		return prov
	}
	// Set up env variables to simulate flags.
	restoreNode := osx.MustSetenv("NODE", "mlab4d.lga0t.measurement-lab.org")
	defer restoreNode()
	t.Run("success", func(t *testing.T) {
		main()
	})

	// main() should successfully add a new node.
	osx.MustSetenv("NODE", "mlab1d.lga0t.measurement-lab.org")
	restoreAdd := osx.MustSetenv("ADD", "1")
	restoreUser := osx.MustSetenv("BMCUSER", "username")
	restorePass := osx.MustSetenv("BMCPASSWORD", "password")
	restoreAddr := osx.MustSetenv("ADDR", "127.0.0.1")
	defer restoreAdd()
	defer restoreUser()
	defer restorePass()
	defer restoreAddr()
	t.Run("success-node-added", func(t *testing.T) {
		main()
	})

	// main() should exit if the node already exists.
	osx.MustSetenv("NODE", "mlab4d.lga0t.measurement-lab.org")
	t.Run("failure-add-node-already-exists", func(t *testing.T) {
		assert.PanicsWithValue(t, "os.Exit called", main,
			"os.Exit was not called")
	})

	osx.MustSetenv("BMCUSER", "")
	// main() should fail when trying to -add without specifying the required
	// arguments (username, password and address).
	t.Run("failure-add-missing-args", func(t *testing.T) {
		assert.PanicsWithValue(t, "os.Exit called", main,
			"os.Exit was not called")
	})
	credsNewProvider = oldCredsNewProvider
}
