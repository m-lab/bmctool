package main

import (
	"context"
	"errors"
	"testing"

	"github.com/m-lab/reboot-service/creds"

	"github.com/m-lab/go/osx"
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
	// No node specified, just print usage and return.
	main()

	// Set up env variables to simulate flags.
	restoreNode := osx.MustSetenv("NODE", "mlab4.lga0t.measurement-lab.org")
	defer restoreNode()

	// Create fake Credentials.
	fakeCreds := &creds.Credentials{
		Address:  "0.0.0.0",
		Hostname: "mlab4.lga0t.measurement-lab.org",
		Username: "username",
		Password: "password",
		Model:    "drac",
	}

	oldCreateProvider := createProvider
	createProvider = func(projectID string, namespace string) creds.Provider {
		return &providerMock{
			returnValue: fakeCreds,
		}
	}
	main()

	// main should return if the Credentials object can't be marshalled.
	marshalJSON = func(interface{}, string, string) ([]byte, error) {
		return nil, errors.New("error while marshalling JSON")
	}
	main()

	// main should return if the Provider returns an error.
	createProvider = func(projectID string, namespace string) creds.Provider {
		return &providerMock{
			returnErr: true,
		}
	}
	main()

	createProvider = oldCreateProvider
}
