package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cac",
		Short: "Cloudentity configuration manager",
	}
	rootConfig = RootConfig{}
)

type RootConfig struct {
	ConfigPath string
	Profile    string
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootConfig.ConfigPath, "config", "", "Path to source configuration file")
	rootCmd.PersistentFlags().StringVar(&rootConfig.Profile, "profile", "", "Configuration profile")

	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(diffCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
