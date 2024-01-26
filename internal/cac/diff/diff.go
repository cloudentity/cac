package diff

import (
	"context"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/go-json-experiment/json"
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
		server1 *models.TreeServer
		server2 *models.TreeServer
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

func Tree(source *models.TreeServer, target *models.TreeServer, opts ...Option) (string, error) {
	var (
		options = &Options{}
		outS    map[string]any
		outT    map[string]any
		bts     []byte
		err     error
	)

	for _, opt := range opts {
		opt(options)
	}

	// marshaling structs to json and back to get proper field names in the comparison
	if bts, err = json.Marshal(source); err != nil {
		return "", err
	}

	if err = json.Unmarshal(bts, &outS); err != nil {
		return "", err
	}

	if bts, err = json.Marshal(target); err != nil {
		return "", err
	}

	if err = json.Unmarshal(bts, &outT); err != nil {
		return "", err
	}

	if options.PresentAtSource {
		for k, _ := range outT {
			if tm, ok := outT[k].(map[string]any); ok {
				OnlyPresentKeys(outS[k], tm)
			}
		}
	}

	var out = cmp.Diff(outT, outS)

	if options.Color {
		return colorize(out), nil
	}

	return out, nil
}
