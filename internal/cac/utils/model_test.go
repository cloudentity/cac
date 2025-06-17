package utils_test

import (
	"reflect"
	"testing"

	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/stretchr/testify/require"
)

func TestFilterPatch(t *testing.T) {
	tcs := []struct {
		name     string
		server   models.Rfc7396PatchOperation
		filters  []string
		expected models.Rfc7396PatchOperation
	}{
		{
			name: "only clients",
			server: models.Rfc7396PatchOperation{
				"clients": models.TreeClients{
					"123": models.TreeClient{
						ClientName: "client1",
					},
				},
				"idps": models.TreeIDPs{
					"456": models.TreeIDP{
						Name: "idp1",
					},
				},
			},
			filters: []string{"clients"},
			expected: models.Rfc7396PatchOperation{
				"clients": models.TreeClients{
					"123": models.TreeClient{
						ClientName: "client1",
					},
				},
			},
		},
		{
			name: "only scopes and ciba",
			server: models.Rfc7396PatchOperation{
				"clients": models.TreeClients{
					"123": models.TreeClient{
						ClientName: "client1",
					},
				},
				"scopes_without_service": models.TreeScopes{
					"456": models.TreeScope{
						Description: "some scope",
					},
				},
				"ciba_authentication_service": models.TreeCIBAAuthenticationService{
					Type: "asd",
				},
			},
			filters: []string{"scopes", "ciba"},
			expected: models.Rfc7396PatchOperation{
				"scopes_without_service": models.TreeScopes{
					"456": models.TreeScope{
						Description: "some scope",
					},
				},
				"ciba_authentication_service": models.TreeCIBAAuthenticationService{
					Type: "asd",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := utils.FilterPatch(tc.server, tc.filters)

			require.NoError(t, err)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
