package storage

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

type MultiStorageConfiguration struct {
	DirPath []string `json:"dir_path"`
}

var DefaultMultiStorageConfig = MultiStorageConfiguration{
	DirPath: []string{"data"},
}

func InitMultiStorage(config MultiStorageConfiguration) (*MultiStorage, error) {
	var storages []Storage

	if len(config.DirPath) == 0 {
		return nil, errors.New("at least one dir_path is required")
	}

	for _, config := range config.DirPath {
		storages = append(storages, InitStorage(Configuration{
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
	Config   MultiStorageConfiguration
}

var _ Storage = &MultiStorage{}

// Store for simplicity stores data in first storage only, it is responsibility of the user to move entities to other storages
func (m *MultiStorage) Store(workspace string, data *models.TreeServer) error {
	return m.Storages[0].Store(workspace, data)
}

// Read data from all storages and merge them
func (m *MultiStorage) Read(workspace string) (models.TreeServer, error) {
	var (
		data models.TreeServer
		err  error
	)

	for i := len(m.Storages) - 1; i >= 0; i-- {
		var data2 models.TreeServer

		if data2, err = m.Storages[i].Read(workspace); err != nil {
			return data, errors.Wrap(err, "failed to read data from storage")
		}

		if err = mergo.Merge(&data, data2, mergo.WithOverride); err != nil {
			return data, errors.Wrap(err, "failed to merge data")
		}

	}

	return data, nil
}
