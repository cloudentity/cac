package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cac",
		Short: "Cloudentity configuration manager",
	}
)

func init() {
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
