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

func Diff(ctx context.Context, source api.Source, target api.Source, workspace string, opts ...Option) (string, error) {
	var (
		server1 models.Rfc7396PatchOperation
		server2 models.Rfc7396PatchOperation
		err     error
	)

	if server1, err = source.Read(ctx, workspace); err != nil {
		return "", err
	}

	if server2, err = target.Read(ctx, workspace); err != nil {
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

	delete(source, "id")
	delete(source, "tenant_id")

	delete(target, "id")
	delete(target, "tenant_id")

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
