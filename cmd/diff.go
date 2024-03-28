package cmd

import (
	"github.com/cloudentity/cac/internal/cac"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/diff"
	"github.com/pkg/errors"
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
				result string
				source api.Source
				target api.Source
				err    error
			)

			slog.
				With("workspace", rootConfig.Workspace).
				With("config", rootConfig.ConfigPath).
				With("profile", rootConfig.Profile).
				With("source", diffConfig.Source).
				With("target", diffConfig.Target).
				With("out", diffConfig.Out).
				Info("Comparing workspace configuration")

			if app, err = cac.InitApp(rootConfig.ConfigPath, rootConfig.Profile, rootConfig.Tenant); err != nil {
				return err
			}

			if source, err = app.PickSource(diffConfig.Source, rootConfig.Tenant); err != nil {
				return err
			}

			if target, err = app.PickSource(diffConfig.Target, rootConfig.Tenant); err != nil {
				return err
			}

			slog.Info("Comparing configurations", "source", source, "target", target)

			if result, err = diff.Diff(cmd.Context(), source, target, rootConfig.Workspace,
				diff.Colorize(diffConfig.Colors),
				diff.OnlyPresent(diffConfig.OnlyPresent),
				diff.Filters(diffConfig.Filters...),
			); err != nil {
				return err
			}

			if diffConfig.Out != "-" {
				if err = os.WriteFile(diffConfig.Out, []byte(result), 0644); err != nil {
					return errors.Wrap(err, "failed to write diff result to file")
				}

				return nil
			}

			if _, err = os.Stdout.Write([]byte(result)); err != nil {
				return errors.Wrap(err, "failed to write diff result to stdout")
			}

			return nil
		},
	}
	diffConfig struct {
		Source      string
		Target      string
		WithSecrets bool
		Colors      bool
		OnlyPresent bool
		Filters     []string
		Out         string
	}
)

func init() {
	diffCmd.PersistentFlags().StringVar(&diffConfig.Source, "source", "", "Source profile name")
	diffCmd.PersistentFlags().StringVar(&diffConfig.Target, "target", "", "Target profile name")
	diffCmd.PersistentFlags().BoolVar(&diffConfig.Colors, "colors", true, "Colorize output")
	diffCmd.PersistentFlags().BoolVar(&diffConfig.OnlyPresent, "only-present", false, "Compare only resources present at source")
	diffCmd.PersistentFlags().StringSliceVar(&diffConfig.Filters, "filter", []string{}, "Compare only selected resources")
	diffCmd.PersistentFlags().StringVar(&diffConfig.Out, "out", "-", "Diff output. It can be a file or '-' for stdout")

	mustMarkRequired(diffCmd, "source", "target")
}
