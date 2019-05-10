package main

import (
	"context"
	"errors"
	"testing"

	"github.com/m-lab/go/osx"
	"github.com/m-lab/reboot-service/creds"
	"github.com/stretchr/testify/assert"
)

type providerMock struct {
	returnValue *creds.Credentials
	returnErr   bool
}

func (p *providerMock) FindCredentials(ctx context.Context, node string) (*creds.Credentials, error) {
	if p.returnErr {
		return nil, errors.New("error while fetching credentials")
	}

	return p.returnValue, nil
}

func Test_main(t *testing.T) {
	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "0.0.0.0",
		Hostname: "mlab4.lga0t.measurement-lab.org",
		Username: "username",
		Password: "password",
		Model:    "drac",
	}

	// Replace osExit and logFatalf so that tests don't stop running.
	osExit = func(code int) {
		if code != 1 {
			t.Fatalf("Expected a 1 exit code, got %d.", code)
		}

		panic("os.Exit called")
	}

	oldCredsNewProvider := credsNewProvider

	t.Run("failure-no-node-provided", func(t *testing.T) {
		// No node specified, just print usage and return.
		assert.PanicsWithValue(t, "os.Exit called", main, "os.Exit was not called")
	})

	// Set up env variables to simulate flags.
	restoreNode := osx.MustSetenv("NODE", "mlab4.lga0t.measurement-lab.org")
	defer restoreNode()

	credsNewProvider = func(projectID string, namespace string) creds.Provider {
		return &providerMock{
			returnValue: fakeCreds,
		}
	}
	t.Run("success", func(t *testing.T) {
		main()
	})

	credsNewProvider = oldCredsNewProvider
}
