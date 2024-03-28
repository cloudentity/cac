package diff

import (
	"context"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/google/go-cmp/cmp"
)

type Options struct {
	Color           bool
	PresentAtSource bool
	Filters         []string
}

type Option func(*Options)

func Colorize(colors bool) Option {
	return func(options *Options) {
		options.Color = colors
	}
}

func OnlyPresent(present bool) Option {
	return func(options *Options) {
		options.PresentAtSource = present
	}
}

func Filters(filters ...string) Option {
	return func(options *Options) {
		options.Filters = filters
	}
}

func Diff(ctx context.Context, source api.Source, target api.Source, workspace string, opts ...Option) (string, error) {
	var (
		server1  models.Rfc7396PatchOperation
		server2  models.Rfc7396PatchOperation
		options  = &Options{}
		readOpts []api.SourceOpt
		err      error
	)

	for _, opt := range opts {
		opt(options)
	}

	if len(options.Filters) > 0 {
		readOpts = append(readOpts, api.WithFilters(options.Filters))
	}

	if workspace != "" {
		readOpts = append(readOpts, api.WithWorkspace(workspace))
	}

	if server1, err = source.Read(ctx, readOpts...); err != nil {
		return "", err
	}

	if server2, err = target.Read(ctx, readOpts...); err != nil {
		return "", err
	}

	return Tree(server1, server2, opts...)
}

func Tree(source models.Rfc7396PatchOperation, target models.Rfc7396PatchOperation, opts ...Option) (string, error) {
	var (
		options = &Options{}
		err     error
	)

	for _, opt := range opts {
		opt(options)
	}

	utils.CleanPatch(source)
	utils.CleanPatch(target)

	// marshaling structs to json and back to get proper field names in the comparison
	if source, err = utils.NormalizePatch(source); err != nil {
		return "", err
	}

	if target, err = utils.NormalizePatch(target); err != nil {
		return "", err
	}

	if options.PresentAtSource {
		for k := range target {
			if tm, ok := target[k].(map[string]any); ok {
				OnlyPresentKeys(source[k], tm)
			}
		}
	}

	var out = cmp.Diff(target, source)

	if options.Color {
		return colorize(out), nil
	}

	return out, nil
}
