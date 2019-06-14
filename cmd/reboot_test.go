package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns a *http.Client with Transport replaced to avoid making
// real calls.
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func Test_reboot(t *testing.T) {
	// Replace osExit so that tests don't stop running.
	oldOsExit := osExit
	osExit = func(code int) {
		if code != 1 {
			t.Fatalf("Expected a 1 exit code, got %d.", code)
		}

		panic("os.Exit called")
	}
	// If REBOOTAPIURL isn't set, reboot() should fail.
	assert.PanicsWithValue(t, "os.Exit called", func() {
		reboot("mlab1d.lga0t.measurement-lab.org")
	}, "os.Exit was not called")
	osExit = oldOsExit

	rebootAPIURL = "dummy"

	// Set up a http.Client that returns values useful for testing.
	oldHTTPClient := httpClient
	httpClient = NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewBufferString(
				"Server power operation successful.")),
			Header: make(http.Header),
		}
	})
	reboot("mlab1d.lga0t.measurement-lab.org")
	httpClient = oldHTTPClient
}
