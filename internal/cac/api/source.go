package api

import (
	"context"
	"errors"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
)

type SourceType string

var ErrUnknownSource = errors.New("unknown source")

const (
	SourceLocal  SourceType = "local"
	SourceRemote SourceType = "remote"
)

func SourceFromString(s string) (SourceType, error) {
	switch s {
	case "local":
		return SourceLocal, nil
	case "remote":
		return SourceRemote, nil
	}

	return "", ErrUnknownSource
}

type Options struct {
	Secrets bool
	Mode    string
	Method  string
}

type SourceOpt func(*Options)

type Source interface {
	Read(ctx context.Context, workspace string, opts ...SourceOpt) (models.Rfc7396PatchOperation, error)
	Write(ctx context.Context, workspace string, data models.Rfc7396PatchOperation, opts ...SourceOpt) error

	String() string
}

func WithSecrets(secrets bool) SourceOpt {
	return func(o *Options) {
		o.Secrets = secrets
	}
}

func WithMode(mode string) SourceOpt {
	return func(o *Options) {
		if mode == "" {
			mode = "update"
		}

		o.Mode = mode
	}
}

func WithMethod(method string) SourceOpt {
	return func(options *Options) {
		options.Method = method
	}
}
