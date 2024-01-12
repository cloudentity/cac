package cmd

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac"
	"github.com/cloudentity/cac/internal/cac/diff"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"os"
)

var (
	diffCmd = &cobra.Command{
		Use:   "diff",
		Short: "Compare configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				app    *cac.Application
				data1  models.TreeServer
				data2  *models.TreeServer
				result any
				bts    []byte
				err    error
			)

			if app, err = cac.InitApp(diffConfig.ConfigPath); err != nil {
				return err
			}

			slog.
				With("workspace", diffConfig.Workspace).
				With("config", diffConfig.ConfigPath).
				Info("Comparing workspace configuration")

			if data1, err = app.Storage.Read(diffConfig.Workspace); err != nil {
				return err
			}

			if data2, err = app.Client.PullWorkspaceConfiguration(cmd.Context(), diffConfig.Workspace, diffConfig.WithSecrets); err != nil {
				return err
			}

			if result, err = diff.Tree(data1, *data2); err != nil {
				return err
			}

			// TODO it should be defined on the result
			bts = result.([]byte)
			_, err = os.Stdout.Write(bts)

			return nil
		},
	}
	diffConfig struct {
		ConfigPath  string
		Workspace   string
		TargetPath  string
		WithSecrets bool
	}
)

func init() {
	diffCmd.PersistentFlags().StringVar(&diffConfig.ConfigPath, "config", "", "Path to configuration file")
	diffCmd.PersistentFlags().StringVar(&diffConfig.Workspace, "workspace", "", "Workspace to load")
	diffCmd.PersistentFlags().BoolVar(&diffConfig.WithSecrets, "with-secrets", false, "Pull secrets")
	diffCmd.PersistentFlags().StringVar(&diffConfig.TargetPath, "target", "remote", "Compare ")
}
