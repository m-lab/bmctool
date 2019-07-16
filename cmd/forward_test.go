package cmd

import "testing"

func Test_splitPorts(t *testing.T) {
	tests := []struct {
		name    string
		ports   string
		srcPort int32
		dstPort int32
		wantErr bool
	}{
		{
			name:    "ok-ports-split",
			ports:   "4443:443",
			srcPort: 4443,
			dstPort: 443,
		},
		{
			name:    "ok-single-port",
			ports:   "443",
			srcPort: 443,
			dstPort: 443,
		},
		{
			name:    "parse-error-src",
			ports:   "foo:443",
			wantErr: true,
		},
		{
			name:    "parse-error-dst",
			ports:   "443:foo",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, dst, err := splitPorts(tt.ports)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitPorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if src != tt.srcPort {
				t.Errorf("splitPorts() got = %v, want %v", src, tt.srcPort)
			}
			if dst != tt.dstPort {
				t.Errorf("splitPorts() got1 = %v, want %v", dst, tt.dstPort)
			}
		})
	}
}
