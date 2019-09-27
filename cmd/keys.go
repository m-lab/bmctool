package cmd

import (
	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage SSH keys on a BMC",
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
