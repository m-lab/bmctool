package cmd

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
	rebootAPIURL = ""
	assert.PanicsWithValue(t, "os.Exit called", func() {
		reboot("mlab1d-lga0t.mlab-sandbox.measurement-lab.org")
	}, "os.Exit was not called")
	rebootAPIURL = "dummy"

	// Set up a http.Client that returns values useful for testing.
	oldHTTPClient := httpClient
	httpClient = NewTestClient(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.String(), "mlab4d") {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: ioutil.NopCloser(bytes.NewBufferString(
					"This is an error.")),
				Header: make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewBufferString(
				"Server power operation successful.")),
			Header: make(http.Header),
		}
	})

	// Rebooting an mlab4 will always fail.
	assert.PanicsWithValue(t, "os.Exit called", func() {
		reboot("mlab4d-lga0t.mlab-sandbox.measurement-lab.org")
	}, "os.Exit was not called")

	// Successful reboot, stdout should contain the expected message.
	var buf bytes.Buffer
	log.SetOutput(&buf)
	reboot("mlab1d-lga0t.mlab-sandbox.measurement-lab.org")
	log.SetOutput(os.Stderr)

	if !strings.Contains(buf.String(), "Server power operation successful.") {
		t.Errorf("Unexpected output: %s", buf.String())
	}

	httpClient = oldHTTPClient
	osExit = oldOsExit

}
