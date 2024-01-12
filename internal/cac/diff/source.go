package diff

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/client"
	"github.com/cloudentity/cac/internal/cac/storage"
)

type Source string

const (
	SourceLocal  Source = "local"
	SourceRemote Source = "remote"
)

type SourceConfig struct {
	Workspace string
	Source    Source
	Storage   storage.Configuration
	Client    client.Configuration
}

func (source SourceConfig) Load() models.TreeServer {

}
