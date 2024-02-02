package storage

import (
	"context"
	"fmt"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

type MultiStorageConfiguration struct {
	DirPath []string `json:"dir_path"`
}

var DefaultMultiStorageConfig = func() *MultiStorageConfiguration {
	return &MultiStorageConfiguration{
		DirPath: []string{"data"},
	}
}

func InitMultiStorage(config *MultiStorageConfiguration) (*MultiStorage, error) {
	var storages []Storage

	if len(config.DirPath) == 0 {
		return nil, errors.New("at least one dir_path is required")
	}

	for _, config := range config.DirPath {
		storages = append(storages, InitStorage(&Configuration{
			DirPath: config,
		}))
	}

	return &MultiStorage{
		Storages: storages,
		Config:   config,
	}, nil
}

type MultiStorage struct {
	Storages []Storage
	Config   *MultiStorageConfiguration
}

var _ Storage = &MultiStorage{}
var _ api.Source = &MultiStorage{}

// Write for simplicity stores data in first storage only, it is responsibility of the user to move entities to other storages
func (m *MultiStorage) Write(ctx context.Context, workspace string, data models.Rfc7396PatchOperation, opts ...api.SourceOpt) error {
	return m.Storages[0].Write(ctx, workspace, data)
}

// Read data from all storages and merge them
func (m *MultiStorage) Read(ctx context.Context, workspace string, opts ...api.SourceOpt) (models.Rfc7396PatchOperation, error) {
	var (
		data = models.Rfc7396PatchOperation{}
		err  error
	)

	for i := len(m.Storages) - 1; i >= 0; i-- {
		var data2 models.Rfc7396PatchOperation

		if data2, err = m.Storages[i].Read(ctx, workspace); err != nil {
			return data, errors.Wrap(err, "failed to read data from storage")
		}

		if err = mergo.Merge(&data, data2, mergo.WithOverride); err != nil {
			return data, errors.Wrap(err, "failed to merge data")
		}

	}

	return data, nil
}

func (m *MultiStorage) String() string {
	return fmt.Sprintf("storage: %v", m.Config.DirPath)
}
