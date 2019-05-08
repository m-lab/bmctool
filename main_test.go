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

	// Replace osExit so that tests don't stop running.
	osExit = func(code int) {
		if code != 1 {
			t.Fatalf("Expected a 1 exit code, got %d.", code)
		}

		panic("os.Exit called")
	}

	// No node specified, just print usage and return.
	assert.PanicsWithValue(t, "os.Exit called", main, "os.Exit was not called")

	// Set up env variables to simulate flags.
	restoreNode := osx.MustSetenv("NODE", "mlab4.lga0t.measurement-lab.org")
	defer restoreNode()

	// main should exit if the Credentials object can't be marshalled.
	oldCredsNewProvider := credsNewProvider
	credsNewProvider = func(projectID string, namespace string) creds.Provider {
		return &providerMock{
			returnValue: fakeCreds,
		}
	}
	oldJSONMarshalIndent := jsonMarshalIndent
	jsonMarshalIndent = func(interface{}, string, string) ([]byte, error) {
		return nil, errors.New("error while marshalling JSON")
	}
	assert.PanicsWithValue(t, "os.Exit called", main, "os.Exit was not called")
	jsonMarshalIndent = oldJSONMarshalIndent

	// main should exit if the Provider returns an error.
	credsNewProvider = func(projectID string, namespace string) creds.Provider {
		return &providerMock{
			returnErr: true,
		}
	}
	assert.PanicsWithValue(t, "os.Exit called", main, "os.Exit was not called")

	credsNewProvider = func(projectID string, namespace string) creds.Provider {
		return &providerMock{
			returnValue: fakeCreds,
		}
	}
	main()

	credsNewProvider = oldCredsNewProvider
}
