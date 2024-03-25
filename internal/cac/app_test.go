package cac_test

import (
	"github.com/cloudentity/cac/internal/cac"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitApp(t *testing.T) {
	t.Run("InitApp", func(t *testing.T) {
		app, err := cac.InitApp("./../../examples/e2e/config.yaml", "", true)
		require.NoError(t, err)

		require.NotNil(t, app)
		require.NotNil(t, app.Config)
		require.NotNil(t, app.Client)
		require.NotNil(t, app.Storage)
	})
}
