package forwarder

import (
	"reflect"
	"testing"
)

func Test_sshForwarder_getPortParams(t *testing.T) {
	tests := []struct {
		name    string
		ports   []Port
		tHost   string
		tUser   string
		dstHost string
		want    []string
	}{
		{
			name: "success",
			ports: []Port{
				{
					Src: 4443,
					Dst: 443,
				},
				{
					Src: 5900,
					Dst: 5900,
				},
			},
			tHost:   "tunnelhost",
			tUser:   "user",
			dstHost: "dsthost",
			want: []string{
				"-L4443:dsthost:443",
				"-L5900:dsthost:5900",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &sshForwarder{
				ports:   tt.ports,
				tHost:   tt.tHost,
				tUser:   tt.tUser,
				dstHost: tt.dstHost,
			}
			if got := f.getPortParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sshForwarder.getPortParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_New(t *testing.T) {
	f := New("tunnelhost", "user", "dsthost", []Port{})
	if f == nil {
		t.Errorf("New() returned nil.")
	}
}
