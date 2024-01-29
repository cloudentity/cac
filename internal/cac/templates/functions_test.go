package templates_test

import (
	"github.com/cloudentity/cac/internal/cac/logging"
	"github.com/cloudentity/cac/internal/cac/templates"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestFunctions(t *testing.T) {
	dir := t.TempDir()
	err := os.Setenv("FOO", "bar")
	require.NoError(t, err)

	err = logging.InitLogging(&logging.Configuration{
		Level: "debug",
	})
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(dir, "included.txt"), []byte("value from included file"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(dir, "dir1", "dir2"), 0755)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		path     string
		template string
		expected string
		err      error
	}{
		{
			name:     "replace existing env variable",
			path:     filepath.Join(dir, "env-test.yaml"),
			template: `key: {{ env "FOO" }}`,
			expected: `key: bar`,
		},
		{
			name:     "fail to replace missing env variable",
			path:     filepath.Join(dir, "env-test.yaml"),
			template: `key: {{ env "OTHER" }}`,
			err:      templates.ErrEnvNotFound,
		},
		{
			name:     "include file from the same directory",
			path:     filepath.Join(dir, "test.yaml"),
			template: `key: {{ include "included.txt" }}`,
			expected: `key: value from included file`,
		},
		{
			name:     "include file from the same directory with dot prefix",
			path:     filepath.Join(dir, "test.yaml"),
			template: `key: {{ include "./included.txt" }}`,
			expected: `key: value from included file`,
		},
		{
			name:     "include file from parent directory",
			path:     filepath.Join(dir, "dir1", "dir2", "test.yaml"),
			template: `key: {{ include "../../included.txt" }}`,
			expected: `key: value from included file`,
		},
		{
			name:     "include file from execution directory",
			path:     filepath.Join(dir, "test.yaml"),
			template: `key: {{ include "/testdata/file.txt" }}`,
			expected: `key: value from testdata/file.txt`,
		},
		{
			name:     "include file from parent of the execution directory",
			path:     filepath.Join(dir, "test.yaml"),
			template: `{{ include "/../../../examples/e2e/vars.yaml" }}`,
			expected: `var1: value1`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := os.WriteFile(tc.path, []byte(tc.template), 0644)
			require.NoError(t, err)

			var outBts []byte
			outBts, err = templates.New(tc.path).Render()

			if tc.err != nil {
				require.Error(t, err)
				require.ErrorAs(t, err, &tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, string(outBts))
			}
		})
	}
}
