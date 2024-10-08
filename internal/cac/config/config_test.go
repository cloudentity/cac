package config_test

import (
	"github.com/cloudentity/cac/internal/cac/config"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestReadingConfiguration(t *testing.T) {
	t.Run("Reads configuration from file", func(t *testing.T) {
		rootConf, err := config.InitConfig("./../../../examples/e2e/config.yaml")
		require.NoError(t, err)

		expectedIssuer, _ := url.Parse("https://postmance.eu.authz.cloudentity.io/postmance/system")

		require.NotNil(t, rootConf)

		conf, err := rootConf.ForProfile("")
		require.NoError(t, err)

		require.NotNil(t, conf.Client)
		require.Equal(t, expectedIssuer, conf.Client.IssuerURL)
		require.Contains(t, conf.Client.Scopes, "manage_configuration")
		require.NotNil(t, conf.Logging)
		require.Equal(t, "info", conf.Logging.Level)
		require.NotNil(t, conf.Storage)
		require.NotEmpty(t, conf.Client.Scopes)
		require.NotEmpty(t, conf.Logging.Level)
		require.NotEmpty(t, conf.Logging.Format)

		profile, err := rootConf.ForProfile("stage")
		require.NoError(t, err)
		require.NotEmpty(t, "aaaaaaaaaaaaa", profile.Client.ClientID)
	})

	t.Run("fail on invalid path", func(t *testing.T) {
		_, err := config.InitConfig("./invalid.json")
		require.Error(t, err)
		require.Equal(t, "open ./invalid.json: no such file or directory", err.Error())
	})

	t.Run("read config from env", func(t *testing.T) {
		t.Setenv("CLIENT_ISSUER_URL", "https://postmance.eu.authz.cloudentity.io/postmance/system")
		t.Setenv("CLIENT_CLIENT_ID", "test-cid1")
		t.Setenv("CLIENT_CLIENT_SECRET", "test-secret")

		// FIXME reading profiles from env variables is not yet supported
		// t.Setenv("PROFILES_STAGE_CLIENT_CLIENT_SECRET", "test-secret")

		rootConf, err := config.InitConfig("")
		require.NoError(t, err)

		conf, err := rootConf.ForProfile("")
		require.NoError(t, err)

		require.NotNil(t, conf)
		require.NotNil(t, conf.Client)
		require.Equal(t, "test-cid1", conf.Client.ClientID)
		require.Equal(t, "test-secret", conf.Client.ClientSecret)
		require.NotNil(t, conf.Client.IssuerURL)
		require.Equal(t, "https://postmance.eu.authz.cloudentity.io/postmance/system", conf.Client.IssuerURL.String())
	})
}
