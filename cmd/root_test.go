package cmd

import "testing"

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
