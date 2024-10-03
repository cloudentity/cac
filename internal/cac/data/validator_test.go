package data_test

import (
	"testing"

	"github.com/cloudentity/cac/internal/cac/data"
	"github.com/stretchr/testify/require"

	"github.com/cloudentity/acp-client-go/clients/hub/models"
)

func TestServerValidator(t *testing.T) {
	validator := &data.TenantValidator{}
	
	t.Run("allow_script_exec_point_deletion", func(t *testing.T) {
		patch := models.Rfc7396PatchOperation{
			"servers": map[string]any{
				"server1": map[string]any{
					"script_execution_points": map[string]any{
						"client_token_minting": map[string]any{
							"cid1": map[string]any{
								"script_id": "",
							},
						},
					},
				},
			},
		}

		err := validator.Validate(&patch)
		require.NoError(t, err)
	})

}