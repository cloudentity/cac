package templates

import (
	"bytes"
	"golang.org/x/exp/slog"
	"os"
	"text/template"
)

type Template struct {
	Path string
}

func New(path string) *Template {
	return &Template{Path: path}
}

func (t *Template) Render() ([]byte, error) {
	var (
		buff = bytes.Buffer{}
		tmpl *template.Template
		bts  []byte
		err  error
	)

	if bts, err = os.ReadFile(t.Path); err != nil {
		return nil, err
	}

	slog.Debug("rendering template", "path", t.Path, "data", string(bts))

	if tmpl, err = template.New(t.Path).Funcs(functions(t)).Parse(string(bts)); err != nil {
		return nil, err
	}

	if err = tmpl.Execute(&buff, t); err != nil {
		return nil, err
	}

	return buff.Bytes(), err
}
