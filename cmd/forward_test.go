package cmd

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/m-lab/bmctool/forwarder"
	"github.com/stretchr/testify/assert"
)

type forwarderMock struct{}

func (fm *forwarderMock) Start(context.Context) error {
	return nil
}

func newForwarderMock(string, string, []forwarder.Port) forwarder.Forwarder {
	return &forwarderMock{}
}

func Test_splitPorts(t *testing.T) {
	tests := []struct {
		name    string
		ports   string
		want    forwarder.Port
		wantErr bool
	}{
		{
			name:  "success-port-pair",
			ports: "4443:443",
			want: forwarder.Port{
				Src: 4443,
				Dst: 443,
			},
		},
		{
			name:  "success-single-port",
			ports: "443",
			want: forwarder.Port{
				Src: 443,
				Dst: 443,
			},
		},
		{
			name:    "failure-invalid-format-whole-string",
			ports:   "invalid",
			wantErr: true,
		},
		{
			name:    "failure-invalid-format-2nd-port",
			ports:   "443:invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitPorts(tt.ports)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitPorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_forward(t *testing.T) {
	// Replace osExit so that tests don't stop running.
	osExit = func(code int) {
		if code != 1 {
			t.Fatalf("Expected a 1 exit code, got %d.", code)
		}

		panic("os.Exit called")
	}

	defer func() {
		osExit = os.Exit
	}()

	// bmctool forward should fail if the tunnel host or the SSH user aren't
	// set.
	assert.PanicsWithValue(t, "os.Exit called", func() { forward("mlab1.tst01") },
		"os.Exit was not called")

	oldNewForwarder := newForwarder
	defer func() {
		newForwarder = oldNewForwarder
	}()
	newForwarder = newForwarderMock
	tunnelHost = "tunnelhost"
	sshUser = "user"
	_, cancel := context.WithCancel(context.Background())

	// forward() only returns when the context is canceled.
	go forward("mlab1.tst01")
	cancel()
}
