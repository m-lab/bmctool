package forwarder

import (
	"testing"
)

func TestSSHForwarder_Start(t *testing.T) {
	// tests := []struct {
	// 	ports   []string
	// 	tHost   string
	// 	dstHost string
	// 	name    string
	// 	wantErr bool
	// }{
	// 	{
	// 		ports:   []string{"443", "5900"},
	// 		tHost:   "eb.measurementlab.net",
	// 		dstHost: "mlab1d.lga0t.measurement-lab.org",
	// 		name:    "test",
	// 	},
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		f := &sshForwarder{
	// 			ports:   tt.ports,
	// 			tHost:   tt.tHost,
	// 			dstHost: tt.dstHost,
	// 		}

	// 		ctx, cancel := context.WithCancel(context.Background())
	// 		err := f.Start(ctx)
	// 		if err != nil {
	// 			t.Errorf("")
	// 		}

	// 		cancel()
	// 	})
	// }
}
