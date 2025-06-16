package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/cloudentity/cac/internal/cac/logging"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/pkg/errors"
)

type Writer[T any] func(name string, it T) error
type FileNameProvider[T any] func(id string, it T) string

func writeFiles[T any](data map[string]T, parent string, fileName FileNameProvider[T]) error {
	var (
		writer Writer[*WithID[T]]
		names  = map[string]int{}
		err    error
	)

	if len(data) == 0 {
		return nil
	}

	if writer, err = YAMLWriter[*WithID[T]](parent); err != nil {
		return err
	}

	for id, it := range data {
		if reflect.ValueOf(it).IsZero() {
			continue
		}

		name := fileName(id, it)
		count, ok := names[name]

		if ok {
			count = count + 1
			name += fmt.Sprintf("-%d", count)
		} else {
			count = 1
		}

		names[name] = count

		if err = writer(name, NewWithID(id, it)); err != nil {
			return err
		}
	}

	return nil
}

func writeFile[T any](data T, path string) error {
	var (
		parent = filepath.Dir(path)
		writer Writer[T]
		err    error
	)

	if reflect.ValueOf(data).IsZero() {
		logging.Trace("skipping empty file", "path", path)
		return nil
	}

	if writer, err = YAMLWriter[T](parent); err != nil {
		return errors.Wrapf(err, "failed to create YAML writer for path %s", parent)
	}

	if err = writer(filepath.Base(path), data); err != nil {
		return errors.Wrapf(err, "failed to write file %s", path)
	}

	return nil
}

func YAMLWriter[T any](dirPath string) (Writer[T], error) {
	var (
		raw Writer[[]byte]
		err error
	)
	if raw, err = RawWriter(dirPath); err != nil {
		return nil, errors.Wrapf(err, "failed to create raw writer for path %s", dirPath)
	}

	return func(name string, it T) error {
		var (
			bts []byte
			err error
		)

		name += ".yaml"

		if bts, err = utils.ToYaml(it); err != nil {
			return errors.Wrapf(err, "failed to marshal %T to yaml", it)
		}

		bts = postProcessMultilineTemplates(bts)

		if err = raw(name, bts); err != nil {
			return errors.Wrapf(err, "failed to write yaml file %s", name)
		}

		return nil
	}, nil
}

func RawWriter(dirPath string) (Writer[[]byte], error) {
	if err := mkDir(dirPath); err != nil {
		return nil, err
	}

	return func(name string, bts []byte) error {
		var (
			file *os.File
			err  error
		)

		logging.Trace("writing file", "path", filepath.Join(dirPath, name), "data", string(bts))

		if name == "" {
			return fmt.Errorf("file name cannot be empty")
		}

		if strings.HasPrefix(name, ".") {
			return fmt.Errorf("file name cannot start with a dot")
		}

		name = normalize(name)

		if file, err = os.Create(filepath.Join(dirPath, name)); err != nil && !os.IsExist(err) {
			return errors.Wrapf(err, "failed to create file %s", filepath.Join(dirPath, name))
		}

		defer file.Close()

		if _, err = file.Write(bts); err != nil {
			return errors.Wrapf(err, "failed to write data to file %s", filepath.Join(dirPath, name))
		}

		logging.Trace("wrote file", "path", filepath.Join(dirPath, name), "data", string(bts))

		return nil
	}, nil
}

func mkDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

var safeFileNameRegexp = regexp.MustCompile(`[\/:*?"<>| ]`)

func normalize(name string) string {
	return safeFileNameRegexp.ReplaceAllString(name, "_")
}

// createMultilineIncludeTemplate creates a template that will be replaced by a multiline include in a post-processing step
// This helps avoid a situation where the yaml library will wrap the value in quotes and escape internal quotes which will break the go template syntax
func createMultilineIncludeTemplate(str string, indent int) string {
	return fmt.Sprintf(`⌘⌘%d include "%s"⌘⌘`, indent, str)
}

var multilineTemplateRegexp = regexp.MustCompile(`⌘⌘(\d+) ([^⌘]+)⌘⌘`)

func postProcessMultilineTemplates(bts []byte) []byte {
	bts = multilineTemplateRegexp.ReplaceAll(bts, []byte("{{ $2 | nindent $1 }}"))

	return bts
}
