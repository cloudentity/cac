package utils_test

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestFilterPatch(t *testing.T) {
	tcs := []struct {
		name     string
		server   models.Rfc7396PatchOperation
		filters  []string
		expected models.Rfc7396PatchOperation
	}{
		{
			name: "filter patch",
			server: models.Rfc7396PatchOperation{
				"clients": map[string]any{
					"123": models.TreeClient{
						ClientName: "client1",
					},
				},
				"idps": map[string]any{
					"456": models.TreeIDP{
						Name: "idp1",
					},
				},
			},
			filters: []string{"clients"},
			expected: models.Rfc7396PatchOperation{
				"clients": map[string]any{
					"123": models.TreeClient{
						ClientName: "client1",
					},
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
