package cmd

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var (
	pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull existing configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				app  *cac.Application
				data *models.TreeServer
				err  error
			)

			if app, err = cac.InitApp(pullConfig.ConfigPath); err != nil {
				return err
			}

			slog.
				With("workspace", pullConfig.Workspace).
				With("config", pullConfig.ConfigPath).
				Info("Pulling workspace configuration")

			if data, err = app.Client.PullWorkspaceConfiguration(cmd.Context(), pullConfig.Workspace, pullConfig.WithSecrets); err != nil {
				return err
			}

			if err = app.Storage.Store(pullConfig.Workspace, data); err != nil {
				return err
			}

			return nil
		},
	}
	pullConfig struct {
		ConfigPath  string
		Workspace   string
		WithSecrets bool
	}
)

func init() {
	pullCmd.PersistentFlags().StringVar(&pullConfig.ConfigPath, "config", "", "Path to configuration file")
	pullCmd.PersistentFlags().StringVar(&pullConfig.Workspace, "workspace", "", "Workspace to load")
	pullCmd.PersistentFlags().BoolVar(&pullConfig.WithSecrets, "with-secrets", false, "Pull secrets")
}
