package storage

import (
	"context"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/pkg/errors"
	"os"
)

type DryStorage struct {
	DelegatedWriter WriterFunc
}

func InitDryStorage(out string, constr Constructor) (*DryStorage, error) {
	var (
		delegatedWriter WriterFunc
		err             error
	)

	if out == "-" {
		delegatedWriter = stdWriter
	} else if out != "" {
		var (
			file *os.File
			info os.FileInfo
		)

		if file, err = os.OpenFile(out, os.O_RDONLY, 0644); err != nil && !os.IsNotExist(err) {
			return nil, err
		} else if err == nil {
			// file already exists
			defer file.Close()

			if info, err = file.Stat(); err != nil {
				return nil, err
			}

			if info.IsDir() {
				dryStorage := constr(&Configuration{
					DirPath: out,
				})

				delegatedWriter = dryStorage.Write
			}
		}

		if delegatedWriter == nil {
			delegatedWriter = flatFileWriter(out)
		}
	} else {
		return nil, errors.New("out cannot be empty")
	}

	return &DryStorage{
		DelegatedWriter: delegatedWriter,
	}, nil
}

type WriterFunc func(ctx context.Context, data models.Rfc7396PatchOperation, opts ...api.SourceOpt) error

func (d *DryStorage) Write(ctx context.Context, data models.Rfc7396PatchOperation, opts ...api.SourceOpt) error {
	return d.DelegatedWriter(ctx, data, opts...)
}

func (d *DryStorage) Read(ctx context.Context, opts ...api.SourceOpt) (models.Rfc7396PatchOperation, error) {
	panic("read operation is not implemented for dry storage")
}

var stdWriter = func(ctx context.Context, data models.Rfc7396PatchOperation, opts ...api.SourceOpt) error {
	var (
		bts []byte
		err error
	)
	if bts, err = utils.ToYaml(data); err != nil {
		return err
	}

	_, err = os.Stdout.Write(bts)
	return err
}

var flatFileWriter = func(out string) WriterFunc {
	return func(ctx context.Context, data models.Rfc7396PatchOperation, opts ...api.SourceOpt) error {
		var (
			bts []byte
			err error
		)

		if bts, err = utils.ToYaml(data); err != nil {
			return err
		}

		if err = os.WriteFile(out, bts, 0644); err != nil {
			return err
		}

		if bts, err = utils.ToYaml(data); err != nil {
			return err
		}

		// file does not exist or is not a directory
		if err = os.WriteFile(out, bts, 0644); err != nil {
			return err
		}

		return nil
	}
}

var _ Storage = &DryStorage{}
