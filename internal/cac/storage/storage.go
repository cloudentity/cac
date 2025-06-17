package storage

import (
	"context"

	"github.com/cloudentity/cac/internal/cac/api"
)

type Storage interface {
	Write(ctx context.Context, data api.Patch, opts ...api.SourceOpt) error
	Read(ctx context.Context, opts ...api.SourceOpt) (api.Patch, error)
}
