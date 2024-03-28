package cmd

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac"
	"github.com/cloudentity/cac/internal/cac/api"
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
				data models.Rfc7396PatchOperation
				err  error
			)

			if app, err = cac.InitApp(rootConfig.ConfigPath, rootConfig.Profile, rootConfig.Tenant); err != nil {
				return err
			}

			slog.
				With("workspace", rootConfig.Workspace).
				With("tenant", rootConfig.Tenant).
				With("filters", pullConfig.Filters).
				With("config", rootConfig.ConfigPath).
				Info("Pulling workspace configuration")

			if data, err = app.Client.Read(
				cmd.Context(),
				api.WithWorkspace(rootConfig.Workspace),
				api.WithSecrets(pullConfig.WithSecrets),
				api.WithFilters(pullConfig.Filters),
			); err != nil {
				return err
			}

			if err = app.Storage.Write(cmd.Context(), data, api.WithWorkspace(rootConfig.Workspace)); err != nil {
				return err
			}

			return nil
		},
	}
	pullConfig struct {
		WithSecrets bool
		Filters     []string
	}
)

func init() {
	pullCmd.PersistentFlags().BoolVar(&pullConfig.WithSecrets, "with-secrets", false, "Pull secrets")
	pullCmd.PersistentFlags().StringSliceVar(&pullConfig.Filters, "filter", []string{}, "Pull only selected resources")
}
