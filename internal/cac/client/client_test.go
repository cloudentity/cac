package client_test

import (
	"context"
	"fmt"
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

		_, err = c.Read(context.Background(), api.WithSecrets(false), api.WithWorkspace("demo"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown client, no client authentication included, or unsupported authentication method")
	})

	t.Run("client pull configuration without filters", func(t *testing.T) {
		testServer := CreateMockServer(t)
		issuer, _ := url.Parse(fmt.Sprintf("%s/demo/system", testServer.URL))

		c, err := client.InitClient(&client.Configuration{
			Insecure: true,
			Config: acpclient.Config{
				IssuerURL:    issuer,
				ClientID:     "fb346c287c4d4e378cbae39aa0c3fe52",
				ClientSecret: "valid_secret",
			},
		})
		require.NoError(t, err)

		data, err := c.Read(
			context.Background(),
			"admin",
			api.WithSecrets(false),
		)

		require.NoError(t, err)

		require.Len(t, data["clients"], 1)
		require.Len(t, data["idps"], 1)
		require.Equal(t, "demo workspace", data["name"])
	})

	t.Run("client pull configuration and filter", func(t *testing.T) {
		testServer := CreateMockServer(t)
		issuer, _ := url.Parse(fmt.Sprintf("%s/demo/system", testServer.URL))

		c, err := client.InitClient(&client.Configuration{
			Insecure: true,
			Config: acpclient.Config{
				IssuerURL:    issuer,
				ClientID:     "fb346c287c4d4e378cbae39aa0c3fe52",
				ClientSecret: "valid_secret",
			},
		})
		require.NoError(t, err)

		data, err := c.Read(
			context.Background(),
			"admin",
			api.WithSecrets(false),
			api.WithFilters([]string{"clients"}),
		)

		require.NoError(t, err)

		require.Len(t, data["clients"], 1)
		require.Nil(t, data["idps"])
	})
}
