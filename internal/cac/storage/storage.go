package storage

import (
	"context"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
)

type Storage interface {
	Write(ctx context.Context, data models.Rfc7396PatchOperation, opts ...api.SourceOpt) error
	Read(ctx context.Context, opts ...api.SourceOpt) (models.Rfc7396PatchOperation, error)
}
