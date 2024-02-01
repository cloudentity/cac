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
}

type SourceOpt func(*Options)

type Source interface {
	Read(ctx context.Context, workspace string, opts ...SourceOpt) (models.Rfc7396PatchOperation, error)
	Write(ctx context.Context, workspace string, data models.Rfc7396PatchOperation) error

	String() string
}
