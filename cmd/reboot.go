package cmd

import (
	"io/ioutil"
	"net/http"

	"github.com/apex/log"
	"github.com/m-lab/go/rtx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const rebootEndpoint = "/v1/reboot"

var (
	rebootAPIURL string

	// rebootCmd represents the reboot command
	rebootCmd = &cobra.Command{
		Use:   "reboot <hostname>",
		Short: "Reboot a BMC using the Reboot API",
		Long: `This command sends a POST request to the Reboot API to reboot the provided node.

The reboot-api-url flag can be also provided via the REBOOTAPIURL environment variable.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			reboot(args[0])
		},
	}

	httpPost = http.Post
)

func init() {
	rootCmd.AddCommand(rebootCmd)

	viper.AutomaticEnv()

	rebootCmd.Flags().StringVar(&rebootAPIURL, "reboot-api-url",
		viper.GetString("REBOOTAPIURL"), "Reboot API URL")
}

func reboot(host string) {
	// Make sure the Reboot API URL has been provided.
	if rebootAPIURL == "" {
		log.Error("The Reboot API URL must be specified (see bmctool help reboot).")
		osExit(1)
	}
	// Make sure the provided host is a valid M-Lab BMC.
	host = makeBMCHostname(host)
	fullURL := rebootAPIURL + rebootEndpoint + "?host=" + host

	log.Infof("POST %s", fullURL)
	resp, err := httpPost(fullURL, "text/plain", nil)
	rtx.Must(err, "Cannot send reboot request")

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	rtx.Must(err, "Cannot read response from %s")

	log.Info(string(body))
}
