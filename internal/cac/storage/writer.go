package storage

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slog"
)

type Writer[T any] func(name string, it T) error
type FileNameProvider[T any] func(id string, it T) string

func writeFiles[T any](data map[string]T, parent string, fileName FileNameProvider[T]) error {
	var (
		writer Writer[*WithID[T]]
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
		slog.Debug("skipping empty file", "path", path)
		return nil
	}

	if writer, err = YAMLWriter[T](parent); err != nil {
		return err
	}

	if err = writer(filepath.Base(path), data); err != nil {
		return err
	}

	return nil
}

func YAMLWriter[T any](dirPath string) (Writer[T], error) {
	var (
		raw Writer[[]byte]
		err error
	)
	if raw, err = RawWriter(dirPath); err != nil {
		return nil, err
	}

	return func(name string, it T) error {
		var (
			m       = make(map[string]any)
			bts     []byte
			decoder *mapstructure.Decoder
			err     error
		)

		name += ".yaml"

		if decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:  &m,
			TagName: "json",
		}); err != nil {
			return err
		}

		if err = decoder.Decode(it); err != nil {
			return err
		}

		if m, err = maputil.Compact(m); err != nil {
			return err
		}

		if bts, err = yaml.Marshal(m); err != nil {
			return err
		}

		bts = postProcessMultilineTemplates(bts)

		if err = raw(name, bts); err != nil {
			return err
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

		slog.Debug("writing file", "path", filepath.Join(dirPath, name), "data", string(bts))

		if name == "" {
			return fmt.Errorf("file name cannot be empty")
		}

		if strings.HasPrefix(name, ".") {
			return fmt.Errorf("file name cannot start with a dot")
		}

		name = normalize(name)

		if file, err = os.Create(filepath.Join(dirPath, name)); err != nil && !os.IsExist(err) {
			return err
		}

		defer file.Close()

		if _, err = file.Write(bts); err != nil {
			return err
		}

		slog.Debug("wrote file", "path", filepath.Join(dirPath, name), "data", string(bts))

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
	matches := multilineTemplateRegexp.FindSubmatch(bts)
	spaces := ""
	indent := 0

	if len(matches) > 1 {
		var (
			indentBts = matches[1]
			err       error
		)

		if indent, err = strconv.Atoi(string(indentBts)); err != nil {
			indent = 0
		}

		spaces = strings.Repeat(" ", indent)
	}

	bts = multilineTemplateRegexp.ReplaceAll(bts, []byte(fmt.Sprintf("|-\n%s{{ $2 | indent $1 }}", spaces)))

	return bts
}
