package cmd

import (
	"testing"
)

func init() {
	nameVersion = "v1"
}

func Test_makeBMCHostname(t *testing.T) {
	tests := []struct {
		name        string
		nameVersion string
		want        string
	}{
		{
			name:        "mlab1d.lga0t",
			nameVersion: "v1",
			want:        "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name:        "mlab1.lga0t",
			nameVersion: "v1",
			want:        "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name:        "mlab1.lga0t.measurement-lab.org",
			nameVersion: "v1",
			want:        "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name:        "mlab1d.lga0t.measurement-lab.org",
			nameVersion: "v1",
			want:        "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name:        "mlab1d.lga0t.blah",
			nameVersion: "v1",
			want:        "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name:        "mlab1-lga0t",
			nameVersion: "v2",
			want:        "mlab1d-lga0t.mlab-sandbox.measurement-lab.org",
		},
		{
			name:        "mlab1d-lga0t",
			nameVersion: "v2",
			want:        "mlab1d-lga0t.mlab-sandbox.measurement-lab.org",
		},
		{
			name:        "mlab4-abc01",
			nameVersion: "v2",
			want:        "mlab4d-abc01.mlab-staging.measurement-lab.org",
		},
		{
			name:        "mlab1-lga0t.lol.example.org",
			nameVersion: "v2",
			want:        "mlab1d-lga0t.mlab-sandbox.measurement-lab.org",
		},
		{
			name:        "mlab4-abc01.lol.example.org",
			nameVersion: "v2",
			want:        "mlab4d-abc01.mlab-staging.measurement-lab.org",
		},
		{
			name:        "mlab1-lga0t.test-project.measurement-lab.org",
			nameVersion: "v2",
			want:        "mlab1d-lga0t.test-project.measurement-lab.org",
		},
		{
			name:        "mlab1d-lga0t.test-project.measurement-lab.org",
			nameVersion: "v2",
			want:        "mlab1d-lga0t.test-project.measurement-lab.org",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectID = ""
			if got := makeBMCHostname(tt.name, tt.nameVersion); got != tt.want {
				t.Errorf("makeBMCHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProjectID(t *testing.T) {
	tests := []struct {
		name string
		host string
		want string
	}{
		{
			name: "prod-node",
			host: "mlab1.abc01.measurement-lab.org",
			want: "mlab-oti",
		},
		{
			name: "staging-node",
			host: "mlab4.abc01.measurement-lab.org",
			want: "mlab-staging",
		},
		{
			name: "sandbox-node",
			host: "mlab4.abc0t.measurement-lab.org",
			want: "mlab-sandbox",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProjectID(tt.host); got != tt.want {
				t.Errorf("getProjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}
