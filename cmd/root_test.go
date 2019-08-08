package cmd

import (
	"testing"
)

func Test_makeBMCHostname(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "mlab1d.lga0t",
			want: "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name: "mlab1.lga0t",
			want: "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name: "mlab1.lga0t.measurement-lab.org",
			want: "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name: "mlab1d.lga0t.measurement-lab.org",
			want: "mlab1d.lga0t.measurement-lab.org",
		},
		{
			name: "mlab1d.lga0t.blah",
			want: "mlab1d.lga0t.measurement-lab.org",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeBMCHostname(tt.name); got != tt.want {
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
