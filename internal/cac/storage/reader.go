package storage

import (
	"github.com/cloudentity/cac/internal/cac/templates"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"os"
	"path/filepath"
)

type ReadFileOpts[T any] struct {
	Constructor func() *T
}
type ReadFileOpt[T any] func(opts *ReadFileOpts[T])

func readFile[T any](path string, out *T, opts ...ReadFileOpt[T]) error {
	var (
		o   = ReadFileOpts[T]{}
		bts []byte
		err error
	)

	for _, opt := range opts {
		opt(&o)
	}

	if filepath.Ext(path) == "" {
		path += ".yaml"
	}

	slog.Debug("reading file", "path", path)

	if bts, err = templates.New(path).Render(); err != nil {
		if os.IsNotExist(err) {
			slog.Debug("file not found", "path", path)
			return nil
		}

		return errors.Wrapf(err, "failed to render template %s", path)
	}

	if out == nil && o.Constructor != nil {
		out = o.Constructor()
	}

	if err = yaml.Unmarshal(bts, out); err != nil {
		return errors.Wrapf(err, "failed to unmarshal template %s", path)
	}

	slog.Debug("read file", "path", path, "data", bts, "out", out)

	return nil
}

func readFiles[M ~map[string]T, T any](path string, out *M) error {
	var (
		dir []os.DirEntry
		err error
	)

	if dir, err = os.ReadDir(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, file := range dir {
		var (
			name = file.Name()
			it   = WithID[T]{}
		)

		if filepath.Ext(name) != ".yaml" {
			slog.Debug("skipping not yaml file", "name", name)
			continue
		}

		if err = readFile(filepath.Join(path, name), &it); err != nil {
			return err
		}

		(*out)[it.ID] = it.Other
	}

	return nil
}
