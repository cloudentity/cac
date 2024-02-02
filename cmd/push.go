package cmd

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac"
	"github.com/cloudentity/cac/internal/cac/storage"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"os"
)

var (
	pushCmd = &cobra.Command{
		Use:   "push",
		Short: "push local configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				app  *cac.Application
				data models.Rfc7396PatchOperation
				serv *models.TreeServer
				err  error
			)

			if app, err = cac.InitApp(rootConfig.ConfigPath, rootConfig.Profile); err != nil {
				return err
			}

			if data, err = app.Storage.Read(cmd.Context(), pushConfig.Workspace); err != nil {
				return err
			}

			if serv, err = utils.FromPatchToTreeServer(data); err != nil {
				return err
			}

			if err = serv.Validate(strfmt.Default); err != nil {
				return err
			}

			if pushConfig.DryRun {
				var bts []byte
				slog.Info("dry run enabled, storing files to disk instead of pushing to server")

				if pushConfig.Out == "-" {
					if bts, err = utils.ToYaml(data); err != nil {
						return err
					}

					_, err = os.Stdout.Write(bts)
					return err
				}

				if pushConfig.Out != "" {
					var (
						file *os.File
						info os.FileInfo
					)

					if file, err = os.OpenFile(pushConfig.Out, os.O_RDONLY, 0644); err != nil && !os.IsNotExist(err) {
						return err
					} else if err == nil {
						// file already exists
						defer file.Close()

						if info, err = file.Stat(); err != nil {
							return err
						}

						if info.IsDir() {
							dryStorage := storage.InitStorage(&storage.Configuration{
								DirPath: pushConfig.Out,
							})

							if err = dryStorage.Write(cmd.Context(), pushConfig.Workspace, data); err != nil {
								return err
							}

							return nil
						}
					}

					if bts, err = utils.ToYaml(data); err != nil {
						return err
					}

					// file does not exist or is not a directory
					if err = os.WriteFile(pushConfig.Out, bts, 0644); err != nil {
						return err
					}

					return nil
				}
			}

			if err = app.Client.Write(cmd.Context(), pushConfig.Workspace, data); err != nil {
				return errors.Wrap(err, "failed to push workspace configuration")
			}

			slog.Info("pushed workspace configuration", "workspace", pushConfig.Workspace)

			return nil
		},
	}
	pushConfig struct {
		Workspace string
		DryRun    bool
		Out       string
	}
)

func init() {
	pushCmd.PersistentFlags().StringVar(&pushConfig.Workspace, "workspace", "", "Workspace to load")
	pushCmd.PersistentFlags().BoolVar(&pushConfig.DryRun, "dry-run", false, "Write files to disk instead of pushing to server")
	pushCmd.PersistentFlags().StringVar(&pushConfig.Out, "out", "-", "Dry execution output. It can be a file, directory or '-' for stdout")
}
