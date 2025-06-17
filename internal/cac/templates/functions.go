package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudentity/cac/internal/cac/logging"
	zb32 "github.com/corvus-ch/zbase32"
	"github.com/pkg/errors"

	"github.com/Masterminds/sprig/v3"
)

func functions(t *Template) template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	funcMap["include"] = include(t)
	funcMap["env"] = env
	funcMap["nindent"] = nindent
	funcMap["zbase32"] = zbase32
	funcMap["apiID"] = apiID
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
		logging.Trace("including file", "path", fp, "data", str)

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
	return "|-\n" + pad + strings.ReplaceAll(v, "\n", "\n"+pad)
}


func zbase32(input string) string {
	return zb32.StdEncoding.EncodeToString([]byte(input))
}

func apiID(serviceID string, method string, path string) string {
	return zbase32(fmt.Sprintf("%s_%s_%s", serviceID, method, path))
}