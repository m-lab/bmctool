package cmd

import (
	"reflect"
	"testing"

	"github.com/m-lab/bmctool/forwarder"
)

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
