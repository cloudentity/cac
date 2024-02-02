package client_test

import (
	"context"
	acpclient "github.com/cloudentity/acp-client-go"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/client"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestClient(t *testing.T) {
	t.Run("Client init success with valid credentials", func(t *testing.T) {
		issuer, _ := url.Parse("https://postmance.eu.authz.cloudentity.io/postmance/system")

		_, err := client.InitClient(&client.Configuration{
			Config: acpclient.Config{
				IssuerURL:    issuer,
				ClientID:     "fb346c287c4d4e378cbae39aa0c3fe52",
				ClientSecret: "-T1siRsUvmE58hB-2I_fWQZW1lLpk_gK76ZziR8Y9QY",
			},
		})

		require.NoError(t, err)
	})

	t.Run("Client init failure when not pointing at valid issuer url", func(t *testing.T) {
		issuer, _ := url.Parse("https://example.com/tid1/aid1")

		_, err := client.InitClient(&client.Configuration{
			Config: acpclient.Config{
				IssuerURL:    issuer,
				ClientID:     "fb346c287c4d4e378cbae39aa0c3fe52",
				ClientSecret: "-T1siRsUvmE58hB-2I_fWQZW1lLpk_gK76ZziR8Y9QY",
			},
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to get well-known endpoints")
	})

	t.Run("client fails on invalid credentials", func(t *testing.T) {
		issuer, _ := url.Parse("https://postmance.eu.authz.cloudentity.io/postmance/system")

		c, err := client.InitClient(&client.Configuration{
			Config: acpclient.Config{
				IssuerURL:    issuer,
				ClientID:     "fb346c287c4d4e378cbae39aa0c3fe52",
				ClientSecret: "invalid_secret",
			},
		})

		require.NoError(t, err)

		_, err = c.Read(context.Background(), "admin", api.WithSecrets(false))
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown client, no client authentication included, or unsupported authentication method")
	})
}
