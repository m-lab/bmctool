package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/m-lab/go/host"
)

type failingReadCloser struct{}

func (*failingReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("Read() error")
}
func (*failingReadCloser) Close() error {
	return errors.New("Close() error")
}

func Test_makeBMCHostname(t *testing.T) {
	fakeNode := host.Name{
		Machine: "mlab1d",
		Site:    "lga0t",
		Project: "mlab-sandbox",
	}

	projectsJSON := `{
		"mlab1-lga0t": "mlab-sandbox"
		}`

	var resp *http.Response

	oldHTTPGet := httpGet
	httpGet = func(string) (*http.Response, error) {
		return resp, nil
	}

	tests := []struct {
		name        string
		nameVersion string
		want        host.Name
	}{
		{
			name: "mlab1-lga0t",
			want: fakeNode,
		},
		{
			name: "mlab1d-lga0t",
			want: fakeNode,
		},
		{
			name: "mlab1-lga0t.lol.example.org",
			want: fakeNode,
		},
		{
			name: "mlab1-lga0t.test-project.measurement-lab.org",
			want: host.Name{
				Project: "test-project",
				Machine: "mlab1d",
				Site:    "lga0t",
			},
		},
		{
			name: "mlab1d-lga0t.test-project.measurement-lab.org",
			want: host.Name{
				Project: "test-project",
				Machine: "mlab1d",
				Site:    "lga0t",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectID = ""
			resp = &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(projectsJSON)),
				StatusCode: http.StatusOK,
			}

			got := makeBMCHostname(tt.name)
			if tt.want.Machine != got.Machine || tt.want.Site != got.Site ||
				tt.want.Project != got.Project {
				t.Errorf("makeBMCHostname() = %v, want %v", got, tt.want)
			}
		})
	}

	httpGet = oldHTTPGet
}

func Test_getProjectID(t *testing.T) {
	node := host.Name{
		Machine: "mlab1",
		Site:    "lga0t",
		Domain:  "measurement-lab.org",
		Version: "v2",
	}

	projectsJSON := `{
		"mlab1-lga0t": "mlab-sandbox"
		}`

	resp := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(projectsJSON)),
		StatusCode: http.StatusOK,
	}

	oldHTTPGet := httpGet
	httpGet = func(string) (*http.Response, error) {
		return resp, nil
	}

	// A project for this site exists - OK case.
	project, err := getProjectID(node)
	if err != nil {
		t.Errorf("getProjectID() returned an error")
	}
	if project != "mlab-sandbox" {
		t.Errorf("getProjectID(): expected %v, got %v", "mlab-sandbox",
			project)
	}

	// A project for this site does not exist.
	node.Machine = "mlab2"
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(projectsJSON))
	project, err = getProjectID(node)
	if err == nil {
		t.Errorf("getProjectID() expected err, got nil")
	}
	if project != "" {
		t.Errorf("getProjectID(): unexpected return value")
	}

	// Reading the response body fails.
	resp = &http.Response{
		Body:       &failingReadCloser{},
		StatusCode: http.StatusOK,
	}
	project, err = getProjectID(node)
	if err == nil {
		t.Errorf("getProjectID(): expected err, got nil")
	}
	if project != "" {
		t.Errorf("getProjectID(): unexpected return value")
	}

	// http.Get returns an error.
	httpGet = func(string) (*http.Response, error) {
		return nil, errors.New("error")
	}
	project, err = getProjectID(node)
	if err == nil {
		t.Errorf("getProjectID(): expected err, got nil")
	}
	if project != "" {
		t.Errorf("getProjectID(): unexpected return value")
	}

	httpGet = oldHTTPGet
}
