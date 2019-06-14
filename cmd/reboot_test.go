package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_reboot(t *testing.T) {
	httpPost = func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewBufferString(
				"Server power operation successful.")),
			Header: make(http.Header),
		}, nil
	}

	reboot("mlab1d.lga0t.measurement-lab.org")
}
