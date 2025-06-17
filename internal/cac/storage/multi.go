package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudentity/cac/internal/cac/api"
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

type Constructor func(config *Configuration) Storage

func InitMultiStorage(config *MultiStorageConfiguration, constr Constructor) (*MultiStorage, error) {
	var storages []Storage

	if len(config.DirPath) == 0 {
		return nil, errors.New("at least one dir_path is required")
	}

	for _, config := range config.DirPath {
		storages = append(storages, constr(&Configuration{
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
func (m *MultiStorage) Write(ctx context.Context, data api.Patch, opts ...api.SourceOpt) error {
	slog.Debug("Writing data to multi storage")
	return m.Storages[0].Write(ctx, data, opts...)
}

// Read data from all storages and merge them
func (m *MultiStorage) Read(ctx context.Context, opts ...api.SourceOpt) (api.Patch, error) {
	var (
		data api.Patch
		err  error
	)

	slog.Debug("Reading data from multi storage")

	for i := len(m.Storages) - 1; i >= 0; i-- {
		var data2 api.Patch

		if data2, err = m.Storages[i].Read(ctx, opts...); err != nil {
			return nil, errors.Wrap(err, "failed to read data from storage")
		}

		if data == nil {
			data = data2
		} else {
			if err = data.Merge(data2); err != nil {
				return nil, errors.Wrap(err, "failed to merge data")
			}
		}

	}

	return data, nil
}

func (m *MultiStorage) String() string {
	return fmt.Sprintf("storage: %v", m.Config.DirPath)
}
