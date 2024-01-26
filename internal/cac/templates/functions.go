package templates

import (
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func functions(t *Template) template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	funcMap["include"] = include(t)
	funcMap["env"] = env
	funcMap["nindent"] = nindent
	return funcMap
}

func include(t *Template) func(string) (string, error) {
	return func(path string) (string, error) {
		var (
			fp  = filepath.Join(filepath.Dir(t.Path), path)
			bts []byte
			str string
			err error
		)

		if strings.HasPrefix(path, "/") {
			fp = path[1:]
		}

		if bts, err = os.ReadFile(fp); err != nil {
			return "", err
		}

		str = string(bts)
		slog.Debug("including file", "path", fp, "data", str)

		return str, nil
	}
}

var ErrEnvNotFound = errors.New("environment variable not found")

func env(key string) (any, error) {
	env := os.Getenv(key)

	if env == "" {
		return nil, errors.Wrapf(ErrEnvNotFound, "environment variable %s not found", key)
	}

	return env, nil
}

func nindent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return "|-\n" + pad + strings.Replace(v, "\n", "\n"+pad, -1)
}
