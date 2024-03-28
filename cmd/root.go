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
	Workspace  string
	Tenant     bool
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootConfig.ConfigPath, "config", "", "Path to source configuration file")
	rootCmd.PersistentFlags().StringVar(&rootConfig.Profile, "profile", "", "Configuration profile")
	rootCmd.PersistentFlags().BoolVar(&rootConfig.Tenant, "tenant", false, "Tenant configuration")
	rootCmd.PersistentFlags().StringVar(&rootConfig.Workspace, "workspace", "", "Workspace configuration")

	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(diffCmd)

	rootCmd.MarkFlagsMutuallyExclusive("workspace", "tenant")
	rootCmd.MarkFlagsOneRequired("workspace", "tenant")
}

func Execute() error {
	return rootCmd.Execute()
}
