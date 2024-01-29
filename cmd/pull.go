package cmd

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac"
	"github.com/cloudentity/cac/internal/cac/client"
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

			if app, err = cac.InitApp(rootConfig.ConfigPath, rootConfig.Profile); err != nil {
				return err
			}

			slog.
				With("workspace", pullConfig.Workspace).
				With("config", rootConfig.ConfigPath).
				Info("Pulling workspace configuration")

			if data, err = app.Client.Read(cmd.Context(), pullConfig.Workspace, client.WithSecrets(pullConfig.WithSecrets)); err != nil {
				return err
			}

			if err = app.Storage.Write(cmd.Context(), pullConfig.Workspace, data); err != nil {
				return err
			}

			return nil
		},
	}
	pullConfig struct {
		Workspace   string
		WithSecrets bool
	}
)

func init() {
	pullCmd.PersistentFlags().StringVar(&pullConfig.Workspace, "workspace", "", "Workspace to load")
	pullCmd.PersistentFlags().BoolVar(&pullConfig.WithSecrets, "with-secrets", false, "Pull secrets")
}
